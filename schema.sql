DROP DATABASE IF EXISTS cookbook;
CREATE DATABASE cookbook CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE cookbook;

CREATE TABLE recipes (
	id INT AUTO_INCREMENT,
	name VARCHAR(1000),
	yields VARCHAR(250),
	updated VARCHAR(100),

	PRIMARY KEY (id)
);

CREATE TABLE ingredients (
	id INT AUTO_INCREMENT,
	recipeid INT,
	value VARCHAR(2500),

	PRIMARY KEY (id),
	FOREIGN KEY (recipeid) REFERENCES recipes (id)
);

CREATE TABLE instructions (
	id INT AUTO_INCREMENT,
	recipeid INT,
	value VARCHAR(20000),

	PRIMARY KEY (id),
	FOREIGN KEY (recipeid) REFERENCES recipes (id)
);
