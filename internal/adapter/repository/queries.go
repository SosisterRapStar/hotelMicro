package repository

const (
	createHotelQuery = `
INSERT INTO hotels (name, city, address)
VALUES ($1, $2, $3)
RETURNING id, name, city, address, created_at;
`

	getHotelByIDQuery = `
SELECT id, name, city, address, created_at
FROM hotels
WHERE id = $1;
`

	listHotelsQuery = `
SELECT id, name, city, address, created_at
FROM hotels
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
`

	createHotelRoomQuery = `
INSERT INTO hotel_rooms (hotel_id, room_type, rooms_total, rooms_available, price)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, hotel_id, room_type, rooms_total, rooms_available, price;
`

	getHotelRoomByIDQuery = `
SELECT id, hotel_id, room_type, rooms_total, rooms_available, price
FROM hotel_rooms
WHERE id = $1;
`

	listHotelRoomsByHotelIDQuery = `
SELECT id, hotel_id, room_type, rooms_total, rooms_available, price
FROM hotel_rooms
WHERE hotel_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;
`

	createHotelBookingQuery = `
INSERT INTO hotel_bookings (user_id, hotel_id, room_id, check_in, check_out, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING booking_id, user_id, hotel_id, room_id, check_in, check_out, status, created_at, updated_at;
`

	getHotelBookingByIDQuery = `
SELECT booking_id, user_id, hotel_id, room_id, check_in, check_out, status, created_at, updated_at
FROM hotel_bookings
WHERE booking_id = $1;
`

	listHotelBookingsQuery = `
SELECT booking_id, user_id, hotel_id, room_id, check_in, check_out, status, created_at, updated_at
FROM hotel_bookings
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
`

	updateHotelBookingStatusQuery = `
UPDATE hotel_bookings
SET status = $1, updated_at = NOW()
WHERE booking_id = $2
RETURNING booking_id, user_id, hotel_id, room_id, check_in, check_out, status, created_at, updated_at;
`
)
