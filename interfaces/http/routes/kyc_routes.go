package routes

import (
	"go-ride-backend/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

// registerKycRoutes mounts /kyc under the driver-protected group.
func registerKycRoutes(driverProtected *gin.RouterGroup, kycHandler *handlers.KycHandler) {
	kycDocuments := driverProtected.Group("/kyc")
	{
		kycDocuments.POST("/documents/upload-url", kycHandler.RequestUploadURL)
		kycDocuments.POST("/documents/confirm", kycHandler.ConfirmUpload)
		kycDocuments.GET("/status", kycHandler.GetStatus)
	}
}
