package vehicle

import (
	"context"
	"errors"
	"testing"

	domainkyc "go-ride-backend/domain/kyc"
	domainvehicle "go-ride-backend/domain/vehicle"

	"github.com/google/uuid"
)

// fakeKycApprovedReader reports every driver as KYC-approved by default,
// so existing activation tests don't need to know about KYC. Set approved
// to false to exercise the gating path.
type fakeKycApprovedReader struct {
	approved bool
}

func newFakeKycApprovedReader() *fakeKycApprovedReader {
	return &fakeKycApprovedReader{approved: true}
}

func (f *fakeKycApprovedReader) IsKycApproved(_ context.Context, _ uuid.UUID) (bool, error) {
	return f.approved, nil
}

// fakeVehicleDocumentApprovalReader reports every vehicle as having an
// approved registration document by default, so existing activation tests
// don't need to know about per-vehicle verification. Set approved to false
// to exercise the gating path.
type fakeVehicleDocumentApprovalReader struct {
	approved bool
}

func newFakeVehicleDocumentApprovalReader() *fakeVehicleDocumentApprovalReader {
	return &fakeVehicleDocumentApprovalReader{approved: true}
}

func (f *fakeVehicleDocumentApprovalReader) IsVehicleDocumentApproved(_ context.Context, _ uuid.UUID) (bool, error) {
	return f.approved, nil
}

func TestActivateUseCaseExecuteDeactivatesPreviousVehicle(t *testing.T) {
	repo := newFakeVehicleRepo()
	registerUC := NewRegisterUseCase(repo)
	driverID := uuid.New()

	first, err := registerUC.Execute(context.Background(), driverID.String(), RegisterVehicleRequest{
		PlateNumber: "A-1", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register first vehicle: %v", err)
	}
	second, err := registerUC.Execute(context.Background(), driverID.String(), RegisterVehicleRequest{
		PlateNumber: "A-2", Color: "Blue", ModelName: "Civic", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register second vehicle: %v", err)
	}

	activateUC := NewActivateUseCase(repo, newFakeKycApprovedReader(), newFakeVehicleDocumentApprovalReader())
	if _, err := activateUC.Execute(context.Background(), driverID.String(), first.ID); err != nil {
		t.Fatalf("activate first vehicle: %v", err)
	}
	resp, err := activateUC.Execute(context.Background(), driverID.String(), second.ID)
	if err != nil {
		t.Fatalf("activate second vehicle: %v", err)
	}
	if !resp.IsActive {
		t.Fatalf("expected second vehicle to be active")
	}

	firstReloaded, err := repo.GetByID(context.Background(), uuid.MustParse(first.ID))
	if err != nil {
		t.Fatalf("reload first vehicle: %v", err)
	}
	if firstReloaded.IsActive {
		t.Fatalf("expected first vehicle to be deactivated after activating second")
	}
}

func TestActivateUseCaseExecuteForbiddenForOtherDriver(t *testing.T) {
	repo := newFakeVehicleRepo()
	registerUC := NewRegisterUseCase(repo)
	ownerID := uuid.New()

	vehicle, err := registerUC.Execute(context.Background(), ownerID.String(), RegisterVehicleRequest{
		PlateNumber: "A-1", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register vehicle: %v", err)
	}

	activateUC := NewActivateUseCase(repo, newFakeKycApprovedReader(), newFakeVehicleDocumentApprovalReader())
	_, err = activateUC.Execute(context.Background(), uuid.New().String(), vehicle.ID)
	if !errors.Is(err, domainvehicle.ErrVehicleForbidden) {
		t.Fatalf("expected ErrVehicleForbidden, got %v", err)
	}
}

func TestActivateUseCaseExecuteBlockedKycNotApproved(t *testing.T) {
	repo := newFakeVehicleRepo()
	registerUC := NewRegisterUseCase(repo)
	driverID := uuid.New()

	vehicle, err := registerUC.Execute(context.Background(), driverID.String(), RegisterVehicleRequest{
		PlateNumber: "A-1", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register vehicle: %v", err)
	}

	activateUC := NewActivateUseCase(repo, &fakeKycApprovedReader{approved: false}, newFakeVehicleDocumentApprovalReader())
	_, err = activateUC.Execute(context.Background(), driverID.String(), vehicle.ID)
	if !errors.Is(err, domainkyc.ErrKycNotApproved) {
		t.Fatalf("expected ErrKycNotApproved, got %v", err)
	}
}

func TestActivateUseCaseExecuteBlockedVehicleNotVerified(t *testing.T) {
	repo := newFakeVehicleRepo()
	registerUC := NewRegisterUseCase(repo)
	driverID := uuid.New()

	vehicle, err := registerUC.Execute(context.Background(), driverID.String(), RegisterVehicleRequest{
		PlateNumber: "A-1", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register vehicle: %v", err)
	}

	// Driver identity is approved, but this specific vehicle's registration
	// document isn't — activation must still be blocked.
	activateUC := NewActivateUseCase(repo, newFakeKycApprovedReader(), &fakeVehicleDocumentApprovalReader{approved: false})
	_, err = activateUC.Execute(context.Background(), driverID.String(), vehicle.ID)
	if !errors.Is(err, domainkyc.ErrVehicleNotVerified) {
		t.Fatalf("expected ErrVehicleNotVerified, got %v", err)
	}
}
