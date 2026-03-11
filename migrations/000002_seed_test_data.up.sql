-- Seed data for hotels service (MySQL)

INSERT INTO hotels (
    id,
    name,
    city,
    address,
    created_at
) VALUES (
    '44444444-4444-4444-4444-444444444444',
    'Test Hotel',
    'Saint Petersburg',
    'Nevsky prospect 1',
    NOW()
)
ON DUPLICATE KEY UPDATE id = id;

INSERT INTO hotel_rooms (
    id,
    hotel_id,
    room_type,
    rooms_total,
    rooms_available,
    price
) VALUES (
    '55555555-5555-5555-5555-555555555555',
    '44444444-4444-4444-4444-444444444444',
    'standard',
    10,
    10,
    2000
)
ON DUPLICATE KEY UPDATE id = id;

