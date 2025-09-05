CREATE TABLE public.users
(
    id       BIGSERIAL      NOT NULL PRIMARY KEY,
    email    VARCHAR UNIQUE NOT NULL,
    password VARCHAR        NOT NULL,
    username     VARCHAR        NOT NULL UNIQUE
);

CREATE TABLE public.dictionary
(
    id      BIGSERIAL NOT NULL PRIMARY KEY,
    arabic  VARCHAR   NOT NULL,
    russian VARCHAR   NOT NULL
);