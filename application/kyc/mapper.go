package kyc

import (
	domainkyc "go-ride-backend/domain/kyc"
)

func ToDocumentResponse(d *domainkyc.DriverDocument) DocumentResponse {
	resp := DocumentResponse{
		ID:              d.ID.String(),
		DocumentType:    d.DocumentType,
		Status:          d.Status,
		RejectionReason: d.RejectionReason,
	}
	if d.VehicleID != nil {
		vID := d.VehicleID.String()
		resp.VehicleID = &vID
	}
	return resp
}
