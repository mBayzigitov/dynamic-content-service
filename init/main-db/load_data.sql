INSERT INTO features(name)
SELECT 'Feature ' || generate_series(1, 500);

INSERT INTO tags(name)
SELECT 'Tag ' || generate_series(1, 1500);