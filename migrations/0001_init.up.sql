-- ========================================
-- UP Migration
-- ========================================
CREATE TABLE public.roles (
                              id SERIAL PRIMARY KEY,
                              code VARCHAR(50) UNIQUE NOT NULL,
                              name VARCHAR(100) NOT NULL,
                              created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                              updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO public.roles (code, name) VALUES
                  ('admin', 'Администратор'),
                  ('moderator', 'Модератор'),
                  ('worker', 'Сотрудник'),
                  ('collector', 'Сборщик'),
                  ('courier', 'Курьер'),
                  ('user', 'Пользователь')
ON CONFLICT (code) DO NOTHING;


CREATE TABLE public.users
(
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    username VARCHAR UNIQUE NOT NULL,
    role_code VARCHAR(50) NOT NULL DEFAULT 'user',


    --Личная информация о пользователе
    first_name VARCHAR NOT NULL DEFAULT '',
    second_name VARCHAR NOT NULL DEFAULT '',
    phone_number VARCHAR NOT NULL DEFAULT '',

    -- Адресная информация
    apartment VARCHAR(50) NOT NULL DEFAULT '',
    house VARCHAR(50) NOT NULL DEFAULT '',
    street VARCHAR(255) NOT NULL DEFAULT '',
    city VARCHAR(100) NOT NULL DEFAULT '',
    region VARCHAR(100) DEFAULT 'Чеченская республика',
    CONSTRAINT fk_role FOREIGN KEY (role_code) REFERENCES roles(code),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.categories
(
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.tags
(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    usage_count INTEGER DEFAULT 0,
    color VARCHAR(7) NOT NULL DEFAULT '#FFFFFF',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE public.catalogs
(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(8,2) NOT NULL,
    amount INT NOT NULL DEFAULT 0,
    discount_percent DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    sku VARCHAR(64) UNIQUE,
    category_id BIGINT NOT NULL,
    image_url TEXT NOT NULL DEFAULT '',
    weight DECIMAL(8,2) NOT NULL,

    -- Рейтинг и отзывы
    rating DECIMAL(3,2) DEFAULT 0.00,
    reviews_count INTEGER DEFAULT 0,


    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_category
        FOREIGN KEY(category_id)
            REFERENCES categories(id)
            ON DELETE RESTRICT
            ON UPDATE CASCADE
);

CREATE TABLE public.catalog_tags
(
    catalog_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    PRIMARY KEY (catalog_id, tag_id),
    CONSTRAINT fk_catalog
        FOREIGN KEY (catalog_id)
            REFERENCES catalogs(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_tag
        FOREIGN KEY (tag_id)
            REFERENCES tags(id)
            ON DELETE CASCADE
);


-- Базовые категории
INSERT INTO public.categories (code, name, description) VALUES
                        ('vegetablesandfruits', 'Овощи и фрукты', 'Свежие и вкусные овощи и фрукты выращенные специально для вас!'),
                        ('drinks', 'Напитки', 'Вода и другие напитки, для утоления жажды любой разновидности'),
                        ('sweets', 'Сладости', 'Сладости на любой вкус и на любой возраст'),
                        ('meatandfish', 'Мясо и рыба', 'Отборное мясо и рыба которое подарит приятное удевление');

-- Базовые теги
INSERT INTO public.tags (name) VALUES
               ('Новинка'),
               ('Хит продаж'),
               ('Скидка'),
               ('Рекомендуем'),
               ('Утоляет жажду');