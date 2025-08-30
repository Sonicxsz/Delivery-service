CREATE TABLE catalog
(
    id       BIGSERIAL  PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(8,2) NOT NULL,
    amount INT NOT NULL DEFAULT 0,
    discount_percent DECIMAL(5, 2) NOT NULL DEFAULT 0.00,
    sku VARCHAR(64)    UNIQUE,
    category_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT category_id
        FOREIGN KEY(category_id)
            REFERENCES categories(id)
            ON DELETE RESTRICT
            ON UPDATE CASCADE
);


CREATE TABLE catalog_tags
(
    catalog_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    PRIMARY KEY (catalog_id, tag_id),
    CONSTRAINT fk_catalog
        FOREIGN KEY (catalog_id)
            REFERENCES categories(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_tag
        FOREIGN KEY (tag_id)
            REFERENCES tags(id)
            ON DELETE CASCADE
);

