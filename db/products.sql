-- Dijalankan otomatis saat API start (idempotent): tabel dibuat jika belum ada,
-- seed hanya diisi kalau tabel masih kosong.
CREATE TABLE IF NOT EXISTS categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255)  NOT NULL,
    brand       VARCHAR(100)  NOT NULL,
    description TEXT          NOT NULL DEFAULT '',
    price       INTEGER       NOT NULL CHECK (price >= 0),
    stock       INTEGER       NOT NULL DEFAULT 0 CHECK (stock >= 0),
    category_id INTEGER       NOT NULL REFERENCES categories(id),
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);

INSERT INTO categories (name) VALUES
    ('Running'),
    ('Sneakers'),
    ('Basketball'),
    ('Football'),
    ('Casual'),
    ('Boots'),
    ('Sandals')
ON CONFLICT (name) DO NOTHING;

INSERT INTO products (name, brand, description, price, stock, category_id)
SELECT v.name, v.brand, v.description, v.price, v.stock, c.id
FROM (VALUES
    (1, 'Nike Pegasus 41', 'Nike', 'Nike Nike Pegasus 41 - sepatu running berkualitas.', 1799000, 10, 'Running'),
    (2, 'Adidas Ultraboost Light', 'Adidas', 'Adidas Adidas Ultraboost Light - sepatu running berkualitas.', 2599000, 17, 'Running'),
    (3, 'ASICS Gel-Kayano 30', 'ASICS', 'ASICS ASICS Gel-Kayano 30 - sepatu running berkualitas.', 2899000, 24, 'Running'),
    (4, 'Hoka Clifton 9', 'Hoka', 'Hoka Hoka Clifton 9 - sepatu running berkualitas.', 2399000, 31, 'Running'),
    (5, 'Brooks Ghost 15', 'Brooks', 'Brooks Brooks Ghost 15 - sepatu running berkualitas.', 2199000, 38, 'Running'),
    (6, 'New Balance Fresh Foam 880', 'New Balance', 'New Balance New Balance Fresh Foam 880 - sepatu running berkualitas.', 1999000, 45, 'Running'),
    (7, 'Saucony Ride 17', 'Saucony', 'Saucony Saucony Ride 17 - sepatu running berkualitas.', 2099000, 11, 'Running'),
    (8, 'Puma Velocity Nitro 3', 'Puma', 'Puma Puma Velocity Nitro 3 - sepatu running berkualitas.', 1899000, 18, 'Running'),
    (9, 'Nike Air Force 1 ''07', 'Nike', 'Nike Nike Air Force 1 ''07 - sepatu sneakers berkualitas.', 1549000, 25, 'Sneakers'),
    (10, 'Adidas Samba OG', 'Adidas', 'Adidas Adidas Samba OG - sepatu sneakers berkualitas.', 1699000, 32, 'Sneakers'),
    (11, 'Adidas Superstar', 'Adidas', 'Adidas Adidas Superstar - sepatu sneakers berkualitas.', 1499000, 39, 'Sneakers'),
    (12, 'Puma Suede Classic', 'Puma', 'Puma Puma Suede Classic - sepatu sneakers berkualitas.', 1299000, 46, 'Sneakers'),
    (13, 'New Balance 550', 'New Balance', 'New Balance New Balance 550 - sepatu sneakers berkualitas.', 1799000, 12, 'Sneakers'),
    (14, 'Reebok Club C 85', 'Reebok', 'Reebok Reebok Club C 85 - sepatu sneakers berkualitas.', 1199000, 19, 'Sneakers'),
    (15, 'Converse Chuck 70 High', 'Converse', 'Converse Converse Chuck 70 High - sepatu sneakers berkualitas.', 1399000, 26, 'Sneakers'),
    (16, 'Nike Dunk Low', 'Nike', 'Nike Nike Dunk Low - sepatu sneakers berkualitas.', 1649000, 33, 'Sneakers'),
    (17, 'Nike LeBron 21', 'Nike', 'Nike Nike LeBron 21 - sepatu basketball berkualitas.', 2799000, 40, 'Basketball'),
    (18, 'Jordan Air Jordan 1 Mid', 'Jordan', 'Jordan Jordan Air Jordan 1 Mid - sepatu basketball berkualitas.', 2199000, 47, 'Basketball'),
    (19, 'Adidas Dame 8', 'Adidas', 'Adidas Adidas Dame 8 - sepatu basketball berkualitas.', 1799000, 13, 'Basketball'),
    (20, 'Under Armour Curry 11', 'Under Armour', 'Under Armour Under Armour Curry 11 - sepatu basketball berkualitas.', 2299000, 20, 'Basketball'),
    (21, 'Puma MB.03', 'Puma', 'Puma Puma MB.03 - sepatu basketball berkualitas.', 2199000, 27, 'Basketball'),
    (22, 'Li-Ning Way of Wade 10', 'Li-Ning', 'Li-Ning Li-Ning Way of Wade 10 - sepatu basketball berkualitas.', 2499000, 34, 'Basketball'),
    (23, 'Nike Kobe 6 Protro', 'Nike', 'Nike Nike Kobe 6 Protro - sepatu basketball berkualitas.', 2599000, 41, 'Basketball'),
    (24, 'Nike Mercurial Vapor 15', 'Nike', 'Nike Nike Mercurial Vapor 15 - sepatu football berkualitas.', 2299000, 48, 'Football'),
    (25, 'Adidas Predator Accuracy', 'Adidas', 'Adidas Adidas Predator Accuracy - sepatu football berkualitas.', 2199000, 14, 'Football'),
    (26, 'Puma Future 7 Ultimate', 'Puma', 'Puma Puma Future 7 Ultimate - sepatu football berkualitas.', 2099000, 21, 'Football'),
    (27, 'Adidas Copa Pure II', 'Adidas', 'Adidas Adidas Copa Pure II - sepatu football berkualitas.', 1999000, 28, 'Football'),
    (28, 'Mizuno Morelia Neo IV', 'Mizuno', 'Mizuno Mizuno Morelia Neo IV - sepatu football berkualitas.', 2399000, 35, 'Football'),
    (29, 'Specs Accelerator Lightspeed', 'Specs', 'Specs Specs Accelerator Lightspeed - sepatu football berkualitas.', 899000, 42, 'Football'),
    (30, 'Ortuseight Forte Pro', 'Ortuseight', 'Ortuseight Ortuseight Forte Pro - sepatu football berkualitas.', 799000, 49, 'Football'),
    (31, 'Vans Old Skool', 'Vans', 'Vans Vans Old Skool - sepatu casual berkualitas.', 999000, 15, 'Casual'),
    (32, 'Vans Sk8-Hi', 'Vans', 'Vans Vans Sk8-Hi - sepatu casual berkualitas.', 1199000, 22, 'Casual'),
    (33, 'Converse Chuck Taylor All Star', 'Converse', 'Converse Converse Chuck Taylor All Star - sepatu casual berkualitas.', 849000, 29, 'Casual'),
    (34, 'Skechers Go Walk 6', 'Skechers', 'Skechers Skechers Go Walk 6 - sepatu casual berkualitas.', 1099000, 36, 'Casual'),
    (35, 'Compass Gazelle Low', 'Compass', 'Compass Compass Gazelle Low - sepatu casual berkualitas.', 429000, 43, 'Casual'),
    (36, 'Ventela Public Low', 'Ventela', 'Ventela Ventela Public Low - sepatu casual berkualitas.', 389000, 50, 'Casual'),
    (37, 'Aerostreet Daily Slip On', 'Aerostreet', 'Aerostreet Aerostreet Daily Slip On - sepatu casual berkualitas.', 279000, 16, 'Casual'),
    (38, 'Crocs Classic Clog', 'Crocs', 'Crocs Crocs Classic Clog - sepatu casual berkualitas.', 649000, 23, 'Casual'),
    (39, 'Timberland 6-Inch Premium', 'Timberland', 'Timberland Timberland 6-Inch Premium - sepatu boots berkualitas.', 3299000, 30, 'Boots'),
    (40, 'Dr. Martens 1460', 'Dr. Martens', 'Dr. Martens Dr. Martens 1460 - sepatu boots berkualitas.', 2999000, 37, 'Boots'),
    (41, 'Caterpillar Colorado', 'Caterpillar', 'Caterpillar Caterpillar Colorado - sepatu boots berkualitas.', 1799000, 44, 'Boots'),
    (42, 'Red Wing Iron Ranger', 'Red Wing', 'Red Wing Red Wing Iron Ranger - sepatu boots berkualitas.', 4999000, 10, 'Boots'),
    (43, 'Wolverine 1000 Mile', 'Wolverine', 'Wolverine Wolverine 1000 Mile - sepatu boots berkualitas.', 4199000, 17, 'Boots'),
    (44, 'Brodo Signore Chelsea', 'Brodo', 'Brodo Brodo Signore Chelsea - sepatu boots berkualitas.', 899000, 24, 'Boots'),
    (45, 'Pedro Chukka Boots', 'Pedro', 'Pedro Pedro Chukka Boots - sepatu boots berkualitas.', 1099000, 31, 'Boots'),
    (46, 'Birkenstock Arizona', 'Birkenstock', 'Birkenstock Birkenstock Arizona - sepatu sandals berkualitas.', 1499000, 38, 'Sandals'),
    (47, 'Teva Hurricane XLT2', 'Teva', 'Teva Teva Hurricane XLT2 - sepatu sandals berkualitas.', 999000, 45, 'Sandals'),
    (48, 'Havaianas Top', 'Havaianas', 'Havaianas Havaianas Top - sepatu sandals berkualitas.', 279000, 11, 'Sandals'),
    (49, 'Nike Victori One Slide', 'Nike', 'Nike Nike Victori One Slide - sepatu sandals berkualitas.', 449000, 18, 'Sandals'),
    (50, 'Adidas Adilette Shower', 'Adidas', 'Adidas Adidas Adilette Shower - sepatu sandals berkualitas.', 349000, 25, 'Sandals')
) AS v(ord, name, brand, description, price, stock, category)
JOIN categories c ON c.name = v.category
WHERE NOT EXISTS (SELECT 1 FROM products)
ORDER BY v.ord;
