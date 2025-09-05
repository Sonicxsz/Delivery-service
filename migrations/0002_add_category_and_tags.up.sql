CREATE TABLE public.categories (
                           id BIGSERIAL PRIMARY KEY,
                           code VARCHAR(255) NOT NULL UNIQUE,
                           name VARCHAR(255) NOT NULL
);

CREATE TABLE public.tags (
                      id BIGSERIAL PRIMARY KEY,
                      name VARCHAR(255) NOT NULL UNIQUE
);