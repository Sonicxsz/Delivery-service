CREATE TABLE users
(
    id       BIGSERIAL      NOT NULL PRIMARY KEY,
    email    VARCHAR UNIQUE NOT NULL,
    password VARCHAR        NOT NULL,
    name     VARCHAR        NOT NULL
);

CREATE TABLE dictionary
(
    id      BIGSERIAL NOT NULL PRIMARY KEY,
    arabic  VARCHAR   NOT NULL,
    russian VARCHAR   NOT NULL
);