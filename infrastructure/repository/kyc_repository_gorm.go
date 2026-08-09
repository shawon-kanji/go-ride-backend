package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainkyc "go-ride-backend/domain/kyc"

	schemamodels "github.com/shawon-kanji/go-ride-db-schema/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KycRepositoryGorm struct {
	db *gorm.DB
}

func NewKycRepositoryGorm(db *gorm.DB) *KycRepositoryGorm {
	return &KycRepositoryGorm{db: db}
}

func toDomainDriverDocument(model *schemamodels.DriverDocument) *domainkyc.DriverDocument {
	return &domainkyc.DriverDocument{
		ID:              model.ID,
		DriverID:        model.DriverID,
		VehicleID:       model.VehicleID,
		DocumentType:    model.DocumentType,
		FileURL:         model.FileURL,
		Status:          model.Status,
		RejectionReason: model.RejectionReason,
		IsCurrent:       model.IsCurrent,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

// Create supersedes any existing current row in the same scope — same
// driverID+documentType for identity documents (vehicle_id IS NULL), same
// vehicleID+documentType for vehicle documents — and inserts the new row as
// current, in one transaction. Scoping the supersede query by vehicle_id
// (rather than just driver_id+document_type) is what lets a driver have
// current registration documents for two different vehicles at once,
// matching the partial unique indexes added in migration 000023.
func (r *KycRepositoryGorm) Create(ctx context.Context, input *domainkyc.CreateInput) (*domainkyc.DriverDocument, error) {
	now := time.Now().UTC()
	model := &schemamodels.DriverDocument{
		ID:           uuid.New(),
		DriverID:     input.DriverID,
		VehicleID:    input.VehicleID,
		DocumentType: input.DocumentType,
		FileURL:      input.FileURL,
		Status:       schemamodels.DriverDocumentStatusUploaded,
		IsCurrent:    true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		supersedeQuery := tx.Model(&schemamodels.DriverDocument{}).
			Where("document_type = ? AND is_current = true", input.DocumentType)
		if input.VehicleID != nil {
			supersedeQuery = supersedeQuery.Where("vehicle_id = ?", *input.VehicleID)
		} else {
			supersedeQuery = supersedeQuery.Where("driver_id = ? AND vehicle_id IS NULL", input.DriverID)
		}
		if err := supersedeQuery.
			Updates(map[string]interface{}{"is_current": false, "updated_at": now}).Error; err != nil {
			return fmt.Errorf("supersede previous document: %w", err)
		}

		if err := tx.Create(model).Error; err != nil {
			return fmt.Errorf("create driver document: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return toDomainDriverDocument(model), nil
}

func (r *KycRepositoryGorm) GetByID(ctx context.Context, id uuid.UUID) (*domainkyc.DriverDocument, error) {
	var model schemamodels.DriverDocument
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainkyc.ErrDocumentNotFound
		}
		return nil, fmt.Errorf("find driver document by id: %w", err)
	}
	return toDomainDriverDocument(&model), nil
}

func (r *KycRepositoryGorm) ListCurrentByDriver(ctx context.Context, driverID uuid.UUID) ([]*domainkyc.DriverDocument, error) {
	var rows []schemamodels.DriverDocument
	if err := r.db.WithContext(ctx).
		Where("driver_id = ? AND is_current = true", driverID).
		Order("document_type asc").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list current driver documents: %w", err)
	}
	docs := make([]*domainkyc.DriverDocument, 0, len(rows))
	for i := range rows {
		docs = append(docs, toDomainDriverDocument(&rows[i]))
	}
	return docs, nil
}

// IsVehicleDocumentApproved reports whether vehicleID has every required
// vehicle-scoped document type (registration, photos, number plate — see
// domainkyc.RequiredVehicleDocumentTypes) as a current, approved row.
func (r *KycRepositoryGorm) IsVehicleDocumentApproved(ctx context.Context, vehicleID uuid.UUID) (bool, error) {
	var approvedTypes []string
	err := r.db.WithContext(ctx).Model(&schemamodels.DriverDocument{}).
		Where("vehicle_id = ? AND is_current = true AND status = ?", vehicleID, schemamodels.DriverDocumentStatusApproved).
		Pluck("document_type", &approvedTypes).Error
	if err != nil {
		return false, fmt.Errorf("check vehicle document approval: %w", err)
	}

	approved := make(map[string]bool, len(approvedTypes))
	for _, t := range approvedTypes {
		approved[t] = true
	}
	for _, t := range domainkyc.RequiredVehicleDocumentTypes {
		if !approved[t] {
			return false, nil
		}
	}
	return true, nil
}
