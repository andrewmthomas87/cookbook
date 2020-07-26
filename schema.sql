CREATE TABLE recipes (
	id SERIAL PRIMARY KEY,
	category TEXT,
	name TEXT,
	yields TEXT,
	updated TEXT
);

CREATE TABLE ingredients (
	id SERIAL PRIMARY KEY,
	recipeid INT REFERENCES recipes (id),
	value TEXT
);

CREATE TABLE instructions (
	id SERIAL PRIMARY KEY,
	recipeid INT REFERENCES recipes (id),
	value TEXT
);
