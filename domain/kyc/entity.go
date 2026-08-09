package kyc

import (
	"time"

	"github.com/google/uuid"
)

const (
	DocumentTypeSelfie              = "selfie"
	DocumentTypeGovtIDFront         = "govt_id_front"
	DocumentTypeGovtIDBack          = "govt_id_back"
	DocumentTypeDrivingLicenseFront = "driving_license_front"
	DocumentTypeDrivingLicenseBack  = "driving_license_back"
	DocumentTypeVehicleRegistration = "vehicle_registration"
	DocumentTypeVehiclePhotoFront   = "vehicle_photo_front"
	DocumentTypeVehiclePhotoBack    = "vehicle_photo_back"
	DocumentTypeVehiclePhotoSide    = "vehicle_photo_side"
	DocumentTypeVehicleNumberPlate  = "vehicle_number_plate"
)

const (
	DocumentStatusUploaded = "uploaded"
	DocumentStatusApproved = "approved"
	DocumentStatusRejected = "rejected"
)

// RequiredIdentityDocumentTypes are the driver-level (not tied to any one
// vehicle) document types a driver must have a current row for before
// drivers.kyc_status can move past not_started to in_review. These prove
// who the driver is, independent of which vehicle they drive.
var RequiredIdentityDocumentTypes = []string{
	DocumentTypeSelfie,
	DocumentTypeGovtIDFront,
	DocumentTypeGovtIDBack,
	DocumentTypeDrivingLicenseFront,
	DocumentTypeDrivingLicenseBack,
}

// RequiredVehicleDocumentTypes are the vehicle-scoped document types a
// vehicle must have a current, approved row for (per vehicle_id) before
// IsVehicleDocumentApproved reports that vehicle as verified. Insurance and
// fitness certificates are deferred per doc/DRIVER_KYC_PLAN.md and would
// join this list as additional vehicle-scoped types.
var RequiredVehicleDocumentTypes = []string{
	DocumentTypeVehicleRegistration,
	DocumentTypeVehiclePhotoFront,
	DocumentTypeVehiclePhotoBack,
	DocumentTypeVehiclePhotoSide,
	DocumentTypeVehicleNumberPlate,
}

// IsIdentityDocumentType reports whether documentType is a driver-level
// identity document (as opposed to a vehicle-scoped document type).
func IsIdentityDocumentType(documentType string) bool {
	for _, t := range RequiredIdentityDocumentTypes {
		if t == documentType {
			return true
		}
	}
	return false
}

// IsVehicleDocumentType reports whether documentType is one of the
// vehicle-scoped document types (requires a vehicle_id on upload).
func IsVehicleDocumentType(documentType string) bool {
	for _, t := range RequiredVehicleDocumentTypes {
		if t == documentType {
			return true
		}
	}
	return false
}

// DriverDocument is a single uploaded KYC document. Identity documents
// belong only to a driver (VehicleID nil); vehicle_registration documents
// are also linked to the specific vehicle they verify.
type DriverDocument struct {
	ID              uuid.UUID
	DriverID        uuid.UUID
	VehicleID       *uuid.UUID
	DocumentType    string
	FileURL         string
	Status          string
	RejectionReason *string
	IsCurrent       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
