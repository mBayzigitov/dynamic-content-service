-- Inserting data into features table
INSERT INTO features(name)
SELECT 'Feature ' || generate_series(1, 500);

-- Inserting data into tags table
INSERT INTO tags(name)
SELECT 'Tag ' || generate_series(1, 2000);

-- Inserting data into banners table
INSERT INTO banners(content, feature_id)
SELECT
    ('{"title": "some_title ' || f.id || '", "description": "Description of Banner ' || f.id || '"}')::jsonb,
    f.id
FROM features f;

-- Inserting data into banners_tags table (assigning random tags to banners)
INSERT INTO banners_tags(banner_id, tag_id)
SELECT
    b.id,
    t.id
FROM banners b
    CROSS JOIN tags t
ORDER BY b.id;