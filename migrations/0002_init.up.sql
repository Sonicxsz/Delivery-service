CREATE TABLE languages (
                           id BIGSERIAL PRIMARY KEY,
                           code INTEGER NOT NULL UNIQUE,
                           name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE parts_of_speech (
                                 id BIGSERIAL PRIMARY KEY,
                                 code INTEGER NOT NULL UNIQUE,
                                 name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE tags (
                      id BIGSERIAL PRIMARY KEY,
                      name VARCHAR(255) NOT NULL UNIQUE
);