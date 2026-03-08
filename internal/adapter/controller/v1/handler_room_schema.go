package v1

type createHotelRoomRequest struct {
	RoomType       string `json:"room_type"`
	RoomsTotal     int    `json:"rooms_total"`
	RoomsAvailable int    `json:"rooms_available"`
	Price          int    `json:"price"`
}
