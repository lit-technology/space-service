CREATE TABLE IF NOT EXISTS space (
	id BIGINT PRIMARY KEY,
	name VARCHAR(90) NOT NULL,
	photo_url VARCHAR(90),
	country SMALLINT NOT NULL,
	state SMALLINT,
	city SMALLINT,
	type SMALLINT NOT NULL,
	user_id BIGINT,
	password VARCHAR(30),
	gender BOOLEAN,
	age_min SMALLINT,
	age_max SMALLINT,
	posts INT NOT NULL DEFAULT 0,
	followers INT NOT NULL DEFAULT 0,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted TIMESTAMP
);

CREATE INDEX IF NOT EXISTS space_name ON space USING BTREE(name);

CREATE TABLE IF NOT EXISTS space_follow (
	space_id BIGINT NOT NULL,
	user_id BIGINT NOT NULL REFERENCES "user"(id),
	PRIMARY KEY(space_id, user_id)
);
