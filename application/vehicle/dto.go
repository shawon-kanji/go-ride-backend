package vehicle

type RegisterVehicleRequest struct {
	PlateNumber string `json:"plate_number" validate:"required,min=2,max=20"`
	Color       string `json:"color" validate:"required,min=2,max=50"`
	ModelName   string `json:"model_name" validate:"required,min=1,max=100"`
	SeatCount   int    `json:"seat_count" validate:"required,min=1,max=20"`
	Category    string `json:"category" validate:"required,oneof=normal luxury"`
}

type UpdateVehicleRequest struct {
	PlateNumber string `json:"plate_number" validate:"required,min=2,max=20"`
	Color       string `json:"color" validate:"required,min=2,max=50"`
	ModelName   string `json:"model_name" validate:"required,min=1,max=100"`
	SeatCount   int    `json:"seat_count" validate:"required,min=1,max=20"`
	Category    string `json:"category" validate:"required,oneof=normal luxury"`
}

type VehicleResponse struct {
	ID          string `json:"id"`
	DriverID    string `json:"driver_id"`
	PlateNumber string `json:"plate_number"`
	Color       string `json:"color"`
	ModelName   string `json:"model_name"`
	SeatCount   int    `json:"seat_count"`
	Category    string `json:"category"`
	IsActive    bool   `json:"is_active"`
}

type ListVehiclesResponse struct {
	Vehicles []VehicleResponse `json:"vehicles"`
}
