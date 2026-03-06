CREATE TABLE hotels (
	id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
	name TEXT NOT NULL,
	city TEXT NOT NULL,
	address TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE hotel_rooms (
	id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
	hotel_id CHAR(36) NOT NULL,
	room_type TEXT NOT NULL,
	rooms_total INT NOT NULL,
	rooms_available INT NOT NULL,
	price INT NOT NULL,
	FOREIGN KEY (hotel_id) REFERENCES hotels(id)
);
CREATE INDEX idx_hotel_rooms_hotel ON hotel_rooms(hotel_id);

CREATE TABLE hotel_bookings (
	booking_id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
	user_id CHAR(36) NOT NULL,
	hotel_id CHAR(36) NOT NULL,
	room_id CHAR(36) NOT NULL,
	check_in DATE NOT NULL,
	check_out DATE NOT NULL,
	status VARCHAR(32) NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (hotel_id) REFERENCES hotels(id),
	FOREIGN KEY (room_id) REFERENCES hotel_rooms(id)
);

CREATE INDEX idx_hotel_bookings_user ON hotel_bookings(user_id);
CREATE INDEX idx_hotel_bookings_hotel ON hotel_bookings(hotel_id);
