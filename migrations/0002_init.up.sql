CREATE TABLE category (
                           id BIGSERIAL PRIMARY KEY,
                           code VARCHAR(255) NOT NULL UNIQUE,
                           name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE tags (
                      id BIGSERIAL PRIMARY KEY,
                      name VARCHAR(255) NOT NULL UNIQUE
);