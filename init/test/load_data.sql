INSERT INTO features(name)
SELECT 'Feature ' || generate_series(1, 10);

INSERT INTO tags(name)
SELECT 'Tag ' || generate_series(1, 20);

INSERT INTO banners(content, feature_id)
SELECT
    ('{"title": "some_title ' || f.id || '", "description": "Description of Banner ' || f.id || '"}')::jsonb,
    f.id
FROM features f;

INSERT INTO banners_tags(banner_id, tag_id)
SELECT
    b.id,
    t.id
FROM banners b
    CROSS JOIN tags t
ORDER BY b.id;