package repository

const (
	createHotelQuery = `
INSERT INTO hotels (id, name, city, address)
VALUES (?, ?, ?, ?);
`

	getHotelByIDQuery = `
SELECT id, name, city, address, created_at
FROM hotels
WHERE id = ?;
`

	listHotelsQuery = `
SELECT id, name, city, address, created_at
FROM hotels
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
`

	createHotelRoomQuery = `
INSERT INTO hotel_rooms (id, hotel_id, room_type, rooms_total, rooms_available, price)
VALUES (?, ?, ?, ?, ?, ?);
`

	getHotelRoomByIDQuery = `
SELECT id, hotel_id, room_type, rooms_total, rooms_available, price
FROM hotel_rooms
WHERE id = ?;
`

	listHotelRoomsByHotelIDQuery = `
SELECT id, hotel_id, room_type, rooms_total, rooms_available, price
FROM hotel_rooms
WHERE hotel_id = ?
ORDER BY id
LIMIT ? OFFSET ?;
`

	createHotelBookingQuery = `
INSERT INTO hotel_bookings (booking_id, user_id, hotel_id, room_id, check_in, check_out, status)
VALUES (?, ?, ?, ?, ?, ?, ?);
`

	getHotelBookingByIDQuery = `
SELECT booking_id, user_id, hotel_id, room_id, check_in, check_out, status, created_at, updated_at
FROM hotel_bookings
WHERE booking_id = ?;
`

	listHotelBookingsQuery = `
SELECT booking_id, user_id, hotel_id, room_id, check_in, check_out, status, created_at, updated_at
FROM hotel_bookings
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
`

	updateHotelBookingStatusQuery = `
UPDATE hotel_bookings
SET status = ?, updated_at = CURRENT_TIMESTAMP
WHERE booking_id = ?;
`
)
