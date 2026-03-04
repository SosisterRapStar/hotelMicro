CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE hotels (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL,
	city TEXT NOT NULL,
	address TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE hotel_rooms (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	hotel_id UUID NOT NULL REFERENCES hotels(id),
	room_type TEXT NOT NULL,
	rooms_total INTEGER NOT NULL,
	rooms_available INTEGER NOT NULL,
	price INTEGER NOT NULL
);
CREATE INDEX idx_hotel_rooms_hotel ON hotel_rooms(hotel_id);

CREATE TABLE hotel_bookings (
	booking_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL,
	hotel_id UUID NOT NULL REFERENCES hotels(id),
	room_id UUID NOT NULL REFERENCES hotel_rooms(id),
	check_in DATE NOT NULL,
	check_out DATE NOT NULL,
	status TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_hotel_bookings_user ON hotel_bookings(user_id);
CREATE INDEX idx_hotel_bookings_hotel ON hotel_bookings(hotel_id);
