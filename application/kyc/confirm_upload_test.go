package kyc

import (
	"context"
	"errors"
	"testing"

	domaindriver "go-ride-backend/domain/driver"
	domainkyc "go-ride-backend/domain/kyc"

	"github.com/google/uuid"
)

// fakeKycRepo mirrors KycRepositoryGorm's scoping rules: Create supersedes
// the current row in the same scope (driver_id+document_type for identity
// documents, vehicle_id+document_type for vehicle documents) rather than
// just driver_id+document_type, so two vehicles can each have their own
// current registration document.
type fakeKycRepo struct {
	docs map[uuid.UUID]*domainkyc.DriverDocument
}

func newFakeKycRepo() *fakeKycRepo {
	return &fakeKycRepo{docs: map[uuid.UUID]*domainkyc.DriverDocument{}}
}

func (r *fakeKycRepo) Create(_ context.Context, input *domainkyc.CreateInput) (*domainkyc.DriverDocument, error) {
	for _, d := range r.docs {
		if !d.IsCurrent || d.DocumentType != input.DocumentType {
			continue
		}
		if input.VehicleID != nil {
			if d.VehicleID != nil && *d.VehicleID == *input.VehicleID {
				d.IsCurrent = false
			}
		} else if d.VehicleID == nil && d.DriverID == input.DriverID {
			d.IsCurrent = false
		}
	}

	doc := &domainkyc.DriverDocument{
		ID:           uuid.New(),
		DriverID:     input.DriverID,
		VehicleID:    input.VehicleID,
		DocumentType: input.DocumentType,
		FileURL:      input.FileURL,
		Status:       domainkyc.DocumentStatusUploaded,
		IsCurrent:    true,
	}
	r.docs[doc.ID] = doc
	return doc, nil
}

func (r *fakeKycRepo) GetByID(_ context.Context, id uuid.UUID) (*domainkyc.DriverDocument, error) {
	if d, ok := r.docs[id]; ok {
		return d, nil
	}
	return nil, domainkyc.ErrDocumentNotFound
}

func (r *fakeKycRepo) ListCurrentByDriver(_ context.Context, driverID uuid.UUID) ([]*domainkyc.DriverDocument, error) {
	var out []*domainkyc.DriverDocument
	for _, d := range r.docs {
		if d.DriverID == driverID && d.IsCurrent {
			out = append(out, d)
		}
	}
	return out, nil
}

// IsVehicleDocumentApproved mirrors KycRepositoryGorm's rule: every
// required vehicle-scoped document type must have a current, approved row.
func (r *fakeKycRepo) IsVehicleDocumentApproved(_ context.Context, vehicleID uuid.UUID) (bool, error) {
	approved := make(map[string]bool)
	for _, d := range r.docs {
		if d.VehicleID != nil && *d.VehicleID == vehicleID && d.IsCurrent && d.Status == domainkyc.DocumentStatusApproved {
			approved[d.DocumentType] = true
		}
	}
	for _, t := range domainkyc.RequiredVehicleDocumentTypes {
		if !approved[t] {
			return false, nil
		}
	}
	return true, nil
}

// approve is a test helper standing in for the manual "UPDATE
// driver_documents SET status='approved'" workflow used until backoffice
// review endpoints exist.
func (r *fakeKycRepo) approve(documentID uuid.UUID) {
	if d, ok := r.docs[documentID]; ok {
		d.Status = domainkyc.DocumentStatusApproved
	}
}

type fakeStorage struct {
	existingKeys map[string]bool
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{existingKeys: map[string]bool{}}
}

func (s *fakeStorage) PresignUploadURL(_ context.Context, key string, _ string) (string, error) {
	return "https://storage.example.com/" + key, nil
}

func (s *fakeStorage) Exists(_ context.Context, key string) (bool, error) {
	return s.existingKeys[key], nil
}

type fakeDriverKycUpdater struct {
	statusByDriver map[uuid.UUID]string
}

func newFakeDriverKycUpdater() *fakeDriverKycUpdater {
	return &fakeDriverKycUpdater{statusByDriver: map[uuid.UUID]string{}}
}

func (f *fakeDriverKycUpdater) UpdateKycStatus(_ context.Context, id uuid.UUID, kycStatus string) (*domaindriver.Driver, error) {
	f.statusByDriver[id] = kycStatus
	return &domaindriver.Driver{ID: id, KycStatus: kycStatus}, nil
}

// fakeVehicleOwnerChecker reports ownership from an explicit map; vehicles
// not present are reported as not owned by anyone.
type fakeVehicleOwnerChecker struct {
	ownerByVehicle map[uuid.UUID]uuid.UUID
}

func newFakeVehicleOwnerChecker() *fakeVehicleOwnerChecker {
	return &fakeVehicleOwnerChecker{ownerByVehicle: map[uuid.UUID]uuid.UUID{}}
}

func (f *fakeVehicleOwnerChecker) own(vehicleID, driverID uuid.UUID) {
	f.ownerByVehicle[vehicleID] = driverID
}

func (f *fakeVehicleOwnerChecker) IsOwnedByDriver(_ context.Context, vehicleID uuid.UUID, driverID uuid.UUID) (bool, error) {
	owner, ok := f.ownerByVehicle[vehicleID]
	return ok && owner == driverID, nil
}

func TestConfirmUploadUseCaseExecuteMissingRequiredDocumentKeepsNotStarted(t *testing.T) {
	repo := newFakeKycRepo()
	storage := newFakeStorage()
	drivers := newFakeDriverKycUpdater()
	uc := NewConfirmUploadUseCase(repo, storage, drivers, newFakeVehicleOwnerChecker())
	driverID := uuid.New()

	key := "drivers/" + driverID.String() + "/selfie/" + uuid.New().String()
	storage.existingKeys[key] = true

	resp, err := uc.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
		DocumentType: domainkyc.DocumentTypeSelfie,
		Key:          key,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != domainkyc.DocumentStatusUploaded {
		t.Fatalf("unexpected status: %+v", resp)
	}
	if _, ok := drivers.statusByDriver[driverID]; ok {
		t.Fatalf("expected kyc_status not to be touched with only one of several required documents uploaded")
	}
}

func TestConfirmUploadUseCaseExecuteAllRequiredIdentityDocumentsMovesToInReview(t *testing.T) {
	repo := newFakeKycRepo()
	storage := newFakeStorage()
	drivers := newFakeDriverKycUpdater()
	uc := NewConfirmUploadUseCase(repo, storage, drivers, newFakeVehicleOwnerChecker())
	driverID := uuid.New()

	for _, docType := range domainkyc.RequiredIdentityDocumentTypes {
		key := "drivers/" + driverID.String() + "/" + docType + "/" + uuid.New().String()
		storage.existingKeys[key] = true
		if _, err := uc.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
			DocumentType: docType,
			Key:          key,
		}); err != nil {
			t.Fatalf("confirm upload for %s: %v", docType, err)
		}
	}

	if drivers.statusByDriver[driverID] != domaindriver.KycStatusInReview {
		t.Fatalf("expected kyc_status in_review once all required identity documents uploaded, got %q", drivers.statusByDriver[driverID])
	}
}

func TestConfirmUploadUseCaseExecuteVehicleRegistrationDoesNotAffectDriverKycStatus(t *testing.T) {
	repo := newFakeKycRepo()
	storage := newFakeStorage()
	drivers := newFakeDriverKycUpdater()
	vehicles := newFakeVehicleOwnerChecker()
	uc := NewConfirmUploadUseCase(repo, storage, drivers, vehicles)
	driverID := uuid.New()
	vehicleID := uuid.New()
	vehicles.own(vehicleID, driverID)

	key := "drivers/" + driverID.String() + "/vehicle_registration/" + vehicleID.String() + "/" + uuid.New().String()
	storage.existingKeys[key] = true

	resp, err := uc.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
		DocumentType: domainkyc.DocumentTypeVehicleRegistration,
		Key:          key,
		VehicleID:    vehicleID.String(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.VehicleID == nil || *resp.VehicleID != vehicleID.String() {
		t.Fatalf("expected document linked to vehicle, got %+v", resp)
	}
	if _, ok := drivers.statusByDriver[driverID]; ok {
		t.Fatalf("expected drivers.kyc_status not to be touched by a vehicle_registration upload")
	}
}

func TestConfirmUploadUseCaseExecuteVehicleRegistrationForbiddenForNonOwnedVehicle(t *testing.T) {
	repo := newFakeKycRepo()
	storage := newFakeStorage()
	drivers := newFakeDriverKycUpdater()
	vehicles := newFakeVehicleOwnerChecker()
	uc := NewConfirmUploadUseCase(repo, storage, drivers, vehicles)
	driverID := uuid.New()
	vehicleID := uuid.New()
	vehicles.own(vehicleID, uuid.New()) // owned by a different driver

	key := "drivers/" + driverID.String() + "/vehicle_registration/" + vehicleID.String() + "/" + uuid.New().String()
	storage.existingKeys[key] = true

	_, err := uc.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
		DocumentType: domainkyc.DocumentTypeVehicleRegistration,
		Key:          key,
		VehicleID:    vehicleID.String(),
	})
	if !errors.Is(err, domainkyc.ErrDocumentForbidden) {
		t.Fatalf("expected ErrDocumentForbidden, got %v", err)
	}
}

func TestConfirmUploadUseCaseExecuteVehicleApprovedOnlyWhenEveryRequiredTypeApproved(t *testing.T) {
	repo := newFakeKycRepo()
	storage := newFakeStorage()
	drivers := newFakeDriverKycUpdater()
	vehicles := newFakeVehicleOwnerChecker()
	uc := NewConfirmUploadUseCase(repo, storage, drivers, vehicles)
	driverID := uuid.New()
	vehicleID := uuid.New()
	vehicles.own(vehicleID, driverID)

	docIDs := map[string]uuid.UUID{}
	for _, docType := range domainkyc.RequiredVehicleDocumentTypes {
		key := "drivers/" + driverID.String() + "/" + docType + "/" + vehicleID.String() + "/" + uuid.New().String()
		storage.existingKeys[key] = true
		resp, err := uc.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
			DocumentType: docType,
			Key:          key,
			VehicleID:    vehicleID.String(),
		})
		if err != nil {
			t.Fatalf("confirm upload for %s: %v", docType, err)
		}
		docIDs[docType] = uuid.MustParse(resp.ID)
	}

	approved, err := repo.IsVehicleDocumentApproved(context.Background(), vehicleID)
	if err != nil {
		t.Fatalf("check approval: %v", err)
	}
	if approved {
		t.Fatalf("expected vehicle not approved before any document is manually approved")
	}

	// Approve all but one required type — still not fully verified.
	for docType, id := range docIDs {
		if docType == domainkyc.DocumentTypeVehicleNumberPlate {
			continue
		}
		repo.approve(id)
	}
	approved, err = repo.IsVehicleDocumentApproved(context.Background(), vehicleID)
	if err != nil {
		t.Fatalf("check approval: %v", err)
	}
	if approved {
		t.Fatalf("expected vehicle not approved while vehicle_number_plate is still unapproved")
	}

	// Approve the last one — now fully verified.
	repo.approve(docIDs[domainkyc.DocumentTypeVehicleNumberPlate])
	approved, err = repo.IsVehicleDocumentApproved(context.Background(), vehicleID)
	if err != nil {
		t.Fatalf("check approval: %v", err)
	}
	if !approved {
		t.Fatalf("expected vehicle approved once every required document type is approved")
	}
}

func TestConfirmUploadUseCaseExecuteTwoVehiclesHaveIndependentCurrentRegistrationDocuments(t *testing.T) {
	repo := newFakeKycRepo()
	storage := newFakeStorage()
	drivers := newFakeDriverKycUpdater()
	vehicles := newFakeVehicleOwnerChecker()
	uc := NewConfirmUploadUseCase(repo, storage, drivers, vehicles)
	driverID := uuid.New()
	vehicleA := uuid.New()
	vehicleB := uuid.New()
	vehicles.own(vehicleA, driverID)
	vehicles.own(vehicleB, driverID)

	for _, v := range []uuid.UUID{vehicleA, vehicleB} {
		key := "drivers/" + driverID.String() + "/vehicle_registration/" + v.String() + "/" + uuid.New().String()
		storage.existingKeys[key] = true
		if _, err := uc.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
			DocumentType: domainkyc.DocumentTypeVehicleRegistration,
			Key:          key,
			VehicleID:    v.String(),
		}); err != nil {
			t.Fatalf("confirm upload for vehicle %s: %v", v, err)
		}
	}

	current, err := repo.ListCurrentByDriver(context.Background(), driverID)
	if err != nil {
		t.Fatalf("list current: %v", err)
	}
	if len(current) != 2 {
		t.Fatalf("expected both vehicles to have an independent current registration document, got %d", len(current))
	}
}

func TestConfirmUploadUseCaseExecuteReuploadSupersedesPreviousDocument(t *testing.T) {
	repo := newFakeKycRepo()
	storage := newFakeStorage()
	drivers := newFakeDriverKycUpdater()
	uc := NewConfirmUploadUseCase(repo, storage, drivers, newFakeVehicleOwnerChecker())
	driverID := uuid.New()

	firstKey := "drivers/" + driverID.String() + "/selfie/" + uuid.New().String()
	storage.existingKeys[firstKey] = true
	first, err := uc.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
		DocumentType: domainkyc.DocumentTypeSelfie,
		Key:          firstKey,
	})
	if err != nil {
		t.Fatalf("first confirm: %v", err)
	}

	secondKey := "drivers/" + driverID.String() + "/selfie/" + uuid.New().String()
	storage.existingKeys[secondKey] = true
	second, err := uc.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
		DocumentType: domainkyc.DocumentTypeSelfie,
		Key:          secondKey,
	})
	if err != nil {
		t.Fatalf("second confirm: %v", err)
	}

	firstReloaded, err := repo.GetByID(context.Background(), uuid.MustParse(first.ID))
	if err != nil {
		t.Fatalf("reload first document: %v", err)
	}
	if firstReloaded.IsCurrent {
		t.Fatalf("expected first document to be superseded after re-upload")
	}

	current, err := repo.ListCurrentByDriver(context.Background(), driverID)
	if err != nil {
		t.Fatalf("list current: %v", err)
	}
	if len(current) != 1 || current[0].ID.String() != second.ID {
		t.Fatalf("expected exactly the re-uploaded document to be current, got %+v", current)
	}
}

func TestConfirmUploadUseCaseExecuteUploadNotFoundInStorage(t *testing.T) {
	repo := newFakeKycRepo()
	storage := newFakeStorage()
	drivers := newFakeDriverKycUpdater()
	uc := NewConfirmUploadUseCase(repo, storage, drivers, newFakeVehicleOwnerChecker())
	driverID := uuid.New()

	key := "drivers/" + driverID.String() + "/selfie/" + uuid.New().String()

	_, err := uc.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
		DocumentType: domainkyc.DocumentTypeSelfie,
		Key:          key,
	})
	if !errors.Is(err, domainkyc.ErrUploadNotConfirmed) {
		t.Fatalf("expected ErrUploadNotConfirmed, got %v", err)
	}
}

func TestConfirmUploadUseCaseExecuteInvalidDocumentType(t *testing.T) {
	repo := newFakeKycRepo()
	storage := newFakeStorage()
	drivers := newFakeDriverKycUpdater()
	uc := NewConfirmUploadUseCase(repo, storage, drivers, newFakeVehicleOwnerChecker())
	driverID := uuid.New()

	_, err := uc.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
		DocumentType: "not_a_real_type",
		Key:          "drivers/" + driverID.String() + "/not_a_real_type/x",
	})
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestConfirmUploadUseCaseExecuteVehicleRegistrationRequiresVehicleID(t *testing.T) {
	repo := newFakeKycRepo()
	storage := newFakeStorage()
	drivers := newFakeDriverKycUpdater()
	uc := NewConfirmUploadUseCase(repo, storage, drivers, newFakeVehicleOwnerChecker())
	driverID := uuid.New()

	_, err := uc.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
		DocumentType: domainkyc.DocumentTypeVehicleRegistration,
		Key:          "drivers/" + driverID.String() + "/vehicle_registration/x",
	})
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}
