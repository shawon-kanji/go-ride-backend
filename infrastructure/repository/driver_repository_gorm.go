package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	domaindriver "go-ride-backend/domain/driver"

	schemamodels "github.com/shawon-kanji/go-ride-db-schema/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DriverRepositoryGorm struct {
	db *gorm.DB
}

func NewDriverRepositoryGorm(db *gorm.DB) *DriverRepositoryGorm {
	return &DriverRepositoryGorm{db: db}
}

func (r *DriverRepositoryGorm) Create(ctx context.Context, driver *domaindriver.Driver) error {
	model := &schemamodels.Driver{
		ID:              driver.ID,
		Email:           driver.Email,
		PasswordHash:    driver.PasswordHash,
		FirstName:       driver.FirstName,
		LastName:        driver.LastName,
		AccountStatus:   driver.AccountStatus,
		IsEmailVerified: driver.IsEmailVerified,
		IsOnline:        driver.IsOnline,
		CreatedAt:       driver.CreatedAt,
		UpdatedAt:       driver.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("create driver: %w", err)
	}
	return nil
}

func toDomainDriver(model *schemamodels.Driver) *domaindriver.Driver {
	return &domaindriver.Driver{
		ID:              model.ID,
		Email:           model.Email,
		PasswordHash:    model.PasswordHash,
		FirstName:       model.FirstName,
		LastName:        model.LastName,
		AccountStatus:   model.AccountStatus,
		IsEmailVerified: model.IsEmailVerified,
		IsOnline:        model.IsOnline,
		IsPaused:        model.IsPaused,
		KycStatus:       model.KycStatus,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func (r *DriverRepositoryGorm) GetByEmail(ctx context.Context, email string) (*domaindriver.Driver, error) {
	var model schemamodels.Driver
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domaindriver.ErrDriverNotFound
		}
		return nil, fmt.Errorf("find driver by email: %w", err)
	}
	return toDomainDriver(&model), nil
}

func (r *DriverRepositoryGorm) GetByID(ctx context.Context, id uuid.UUID) (*domaindriver.Driver, error) {
	var model schemamodels.Driver
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domaindriver.ErrDriverNotFound
		}
		return nil, fmt.Errorf("find driver by id: %w", err)
	}
	return toDomainDriver(&model), nil
}

func (r *DriverRepositoryGorm) UpdateProfile(ctx context.Context, id uuid.UUID, firstName string, lastName string) (*domaindriver.Driver, error) {
	updates := map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
		"updated_at": time.Now().UTC(),
	}

	result := r.db.WithContext(ctx).Model(&schemamodels.Driver{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("update driver profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, domaindriver.ErrDriverNotFound
	}

	return r.GetByID(ctx, id)
}

// UpdateOnlineStatus flips drivers.is_online and, in the same transaction,
// opens or closes a driver_online_sessions row -- the audit trail
// POLISHED-MVP.md's C4 needed, since is_online alone can't answer "how long
// has this driver been online today". The session write is idempotent
// either direction: going online is an INSERT ... ON CONFLICT (driver_id)
// WHERE ended_at IS NULL DO NOTHING (a driver already online again is a
// no-op, not a duplicate open session — see the matching partial unique
// index in go-ride-db-schema), and closing when no session is open is
// simply a zero-row UPDATE, not an error (the primary is_online mutation's
// own RowsAffected check below is what still guards "driver not found").
func (r *DriverRepositoryGorm) UpdateOnlineStatus(ctx context.Context, id uuid.UUID, isOnline bool) (*domaindriver.Driver, error) {
	now := time.Now().UTC()

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"is_online":  isOnline,
			"updated_at": now,
		}
		// Going offline always clears a stale pause in the same statement:
		// chk_drivers_is_paused_requires_online would otherwise reject this
		// update outright for a driver who was paused when their shift ended.
		if !isOnline {
			updates["is_paused"] = false
		}
		result := tx.Model(&schemamodels.Driver{}).Where("id = ?", id).Updates(updates)
		if result.Error != nil {
			return fmt.Errorf("update driver online status: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return domaindriver.ErrDriverNotFound
		}

		if isOnline {
			if err := tx.Exec(
				`INSERT INTO driver_online_sessions (id, driver_id, started_at, created_at, updated_at)
				 VALUES (gen_random_uuid(), ?, ?, ?, ?)
				 ON CONFLICT (driver_id) WHERE ended_at IS NULL DO NOTHING`,
				id, now, now, now,
			).Error; err != nil {
				return fmt.Errorf("open driver online session: %w", err)
			}
			return nil
		}

		if err := tx.Exec(
			`UPDATE driver_online_sessions SET ended_at = ?, updated_at = ? WHERE driver_id = ? AND ended_at IS NULL`,
			now, now, id,
		).Error; err != nil {
			return fmt.Errorf("close driver online session: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *DriverRepositoryGorm) UpdatePauseStatus(ctx context.Context, id uuid.UUID, isPaused bool) (*domaindriver.Driver, error) {
	updates := map[string]interface{}{
		"is_paused":  isPaused,
		"updated_at": time.Now().UTC(),
	}

	result := r.db.WithContext(ctx).Model(&schemamodels.Driver{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("update driver pause status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, domaindriver.ErrDriverNotFound
	}

	return r.GetByID(ctx, id)
}

// IsKycApproved satisfies application/vehicle's narrow DriverKycStatusReader
// port, so that package can gate vehicle activation on KYC approval without
// importing domain/driver's full Driver type or Repository interface.
func (r *DriverRepositoryGorm) IsKycApproved(ctx context.Context, driverID uuid.UUID) (bool, error) {
	driver, err := r.GetByID(ctx, driverID)
	if err != nil {
		return false, err
	}
	return driver.KycStatus == domaindriver.KycStatusApproved, nil
}

func (r *DriverRepositoryGorm) UpdateKycStatus(ctx context.Context, id uuid.UUID, kycStatus string) (*domaindriver.Driver, error) {
	updates := map[string]interface{}{
		"kyc_status": kycStatus,
		"updated_at": time.Now().UTC(),
	}

	result := r.db.WithContext(ctx).Model(&schemamodels.Driver{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("update driver kyc status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, domaindriver.ErrDriverNotFound
	}

	return r.GetByID(ctx, id)
}
