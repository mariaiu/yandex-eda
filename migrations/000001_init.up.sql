CREATE TABLE IF NOT EXISTS restaurant (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    slug VARCHAR NOT NULL UNIQUE ,
    minimal_delivery_cost DECIMAL,
    rating DECIMAL
);

CREATE TABLE IF NOT EXISTS position (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER,
    name VARCHAR NOT NULL,
    price DECIMAL NOT NULL,
    description VARCHAR NOT NULL,
    weight INTEGER,
    date_of_parsing DATE NOT NULL DEFAULT CURRENT_DATE,
    CONSTRAINT restaurant_fk FOREIGN KEY (restaurant_id) REFERENCES restaurant(id) ON DELETE CASCADE
);
