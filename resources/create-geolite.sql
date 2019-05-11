CREATE TABLE IF NOT EXISTS geolite (
	id INT PRIMARY KEY,
	space_id BIGINT,
	continent_code CHARACTER(2),
	continent_name VARCHAR (13),
	country_code CHARACTER(2),
	country_name VARCHAR(44),
	state_code CHARACTER(3),
	state_name VARCHAR(52),
	area_code CHARACTER(3),
	area_name VARCHAR(38),
	city_name VARCHAR(63)
);

CREATE INDEX IF NOT EXISTS geolite_country_name ON geolite USING HASH(country_name) WHERE state_name IS NULL AND city_name IS NULL;

CREATE INDEX IF NOT EXISTS geolite_state_name ON geolite USING HASH(state_name) WHERE city_name IS NULL;

CREATE INDEX IF NOT EXISTS geolite_city_name ON geolite USING HASH(city_name);

CREATE TABLE IF NOT EXISTS ip (
	ip INET NOT NULL,
	geolite_id INT,
	cell_id BIGINT NOT NULL,
	space_id BIGINT
);

CREATE INDEX IF NOT EXISTS ip_idx ON ip USING HASH(ip);

CREATE TABLE IF NOT EXISTS place (
	space_id BIGINT PRIMARY KEY,
	country_name VARCHAR(44),
	state_name VARCHAR(52),
	city_name VARCHAR(63),
	street_name VARCHAR(100),
	place_name VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS place_street_name ON place USING HASH(street_name)
	WHERE place_name IS NULL;

CREATE INDEX IF NOT EXISTS place_place_name ON place USING HASH(place_name);
