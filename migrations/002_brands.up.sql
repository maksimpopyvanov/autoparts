CREATE TABLE IF NOT EXISTS brands (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS part_brands (
    part_id  INTEGER REFERENCES parts(id) ON DELETE CASCADE,
    brand_id INTEGER REFERENCES brands(id) ON DELETE CASCADE,
    PRIMARY KEY (part_id, brand_id)
);

INSERT INTO brands (name) VALUES
    ('Lada'),
    ('Toyota'),
    ('BMW'),
    ('Mercedes-Benz'),
    ('Volkswagen'),
    ('Kia'),
    ('Hyundai'),
    ('Ford'),
    ('Renault'),
    ('Nissan')
ON CONFLICT (name) DO NOTHING;
