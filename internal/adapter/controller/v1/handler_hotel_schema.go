package v1

type createHotelRequest struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	Address string `json:"address"`
}
