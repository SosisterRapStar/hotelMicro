package v1

type createHotelBookingRequest struct {
	UserID   string `json:"user_id"`
	HotelID  string `json:"hotel_id"`
	RoomID   string `json:"room_id"`
	CheckIn  string `json:"check_in"`  // YYYY-MM-DD
	CheckOut string `json:"check_out"`  // YYYY-MM-DD
	Status   string `json:"status,omitempty"`
}
