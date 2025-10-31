INSERT INTO users (name)
SELECT 'Demo'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE name = 'Demo');

INSERT INTO users (name)
SELECT 'Other'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE name = 'Other');

INSERT INTO categories (title, is_global)
SELECT 'Food', TRUE
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE is_global = TRUE AND title = 'Food');

INSERT INTO categories (title, is_global)
SELECT 'Transport', TRUE
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE is_global = TRUE AND title = 'Transport');

INSERT INTO categories (title, is_global)
SELECT 'Health', TRUE
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE is_global = TRUE AND title = 'Health');
