package saga

const (
	insertHotelBookingQuery = `
INSERT INTO hotel_bookings (booking_id, user_id, hotel_id, room_id, check_in, check_out, status)
VALUES (?, ?, ?, ?, ?, ?, 'reserved');
`

	decrementRoomAvailableQuery = `
UPDATE hotel_rooms
SET rooms_available = rooms_available - 1
WHERE id = ? AND rooms_available > 0;
`

	cancelHotelBookingQuery = `
UPDATE hotel_bookings
SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
WHERE booking_id = ? AND status = 'reserved';
`

	incrementRoomAvailableQuery = `
UPDATE hotel_rooms
SET rooms_available = rooms_available + 1
WHERE id = (SELECT room_id FROM hotel_bookings WHERE booking_id = ? LIMIT 1);
`
)
