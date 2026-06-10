CREATE TABLE IF NOT EXISTS categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS parts (
    id          SERIAL PRIMARY KEY,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    name        VARCHAR(200) NOT NULL,
    article     VARCHAR(100) UNIQUE,
    description TEXT,
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stock (
    id       SERIAL PRIMARY KEY,
    part_id  INTEGER REFERENCES parts(id) ON DELETE CASCADE UNIQUE,
    quantity INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS income (
    id       SERIAL PRIMARY KEY,
    part_id  INTEGER REFERENCES parts(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    date     DATE NOT NULL DEFAULT CURRENT_DATE,
    comment  TEXT
);

CREATE TABLE IF NOT EXISTS outcome (
    id       SERIAL PRIMARY KEY,
    part_id  INTEGER REFERENCES parts(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    date     DATE NOT NULL DEFAULT CURRENT_DATE,
    comment  TEXT
);
