DROP TABLE IF EXISTS features;
CREATE TABLE IF NOT EXISTS features(
    id   INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name VARCHAR(255)
);

DROP TABLE IF EXISTS tags;
CREATE TABLE IF NOT EXISTS tags(
    id   INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name VARCHAR(255)
);

DROP TABLE IF EXISTS banners;
CREATE TABLE IF NOT EXISTS banners(
    id         INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    content    JSONB,
    feature_id INT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at timestamp DEFAULT now(),
    is_active  BOOL DEFAULT true,
    FOREIGN KEY (feature_id) REFERENCES features (id)
);

DROP TABLE IF EXISTS banners_tags;
CREATE TABLE IF NOT EXISTS banners_tags(
    banner_id INT,
    tag_id    INT,
    FOREIGN KEY (banner_id) REFERENCES banners (id),
    FOREIGN KEY (tag_id) REFERENCES tags (id),
    CONSTRAINT banners_tags_pk PRIMARY KEY(banner_id, tag_id)
);