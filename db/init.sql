SET NAMES utf8mb4;
SET character_set_client = utf8mb4;
SET character_set_connection = utf8mb4;
SET character_set_results = utf8mb4;

-- Create database
CREATE DATABASE IF NOT EXISTS convenience_store CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE convenience_store;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(64) PRIMARY KEY,
    wechat_open_id VARCHAR(128) NOT NULL UNIQUE,
    nickname VARCHAR(128) NOT NULL,
    avatar_url VARCHAR(255) NOT NULL,
    phone VARCHAR(32) DEFAULT NULL,
    default_address_id VARCHAR(64) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Addresses table
CREATE TABLE IF NOT EXISTS addresses (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    recipient VARCHAR(128) NOT NULL,
    phone VARCHAR(32) NOT NULL,
    province VARCHAR(64) NOT NULL,
    city VARCHAR(64) NOT NULL,
    district VARCHAR(64) NOT NULL,
    detail VARCHAR(255) NOT NULL,
    postal_code VARCHAR(16) DEFAULT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_addresses_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Products table
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    tags JSON NULL,
    images JSON NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Cart items table
CREATE TABLE IF NOT EXISTS cart_items (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    product_id VARCHAR(64) NOT NULL,
    quantity INT NOT NULL,
    selected BOOLEAN NOT NULL DEFAULT TRUE,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_cart_items_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_cart_items_products FOREIGN KEY (product_id) REFERENCES products(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    address_id VARCHAR(64) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_orders_users FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_orders_addresses FOREIGN KEY (address_id) REFERENCES addresses(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Order items table
CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL,
    product_id VARCHAR(64) NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    CONSTRAINT fk_order_items_orders FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_order_items_products FOREIGN KEY (product_id) REFERENCES products(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Seed products
INSERT INTO products (id, name, description, price, stock, tags, images, is_active)
VALUES
    ('sku_energy', 'Energy Drink', 'Energy drink to boost performance', 4.50, 500, JSON_ARRAY('drink', 'energy'), JSON_ARRAY('/images/products/energy-1.png', '/images/products/energy-2.png'), TRUE),
    ('sku_snack', 'Potato Chips', 'Classic potato chips', 6.80, 300, JSON_ARRAY('snack'), JSON_ARRAY('/images/products/chips-1.png'), TRUE),
    ('sku_noodle', 'Instant Ramen', 'Quick instant ramen', 8.50, 200, JSON_ARRAY('instant', 'noodle'), JSON_ARRAY('/images/products/ramen-1.png', '/images/products/ramen-2.png'), TRUE)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    description = VALUES(description),
    price = VALUES(price),
    stock = VALUES(stock),
    tags = VALUES(tags),
    images = VALUES(images),
    is_active = VALUES(is_active);

-- Seed user and address
INSERT INTO users (id, wechat_open_id, nickname, avatar_url, phone)
VALUES
    ('usr_demo', 'wx_openid_demo', 'Demo User', 'https://example.com/avatar.png', '18800000000')
ON DUPLICATE KEY UPDATE
    nickname = VALUES(nickname),
    avatar_url = VALUES(avatar_url),
    phone = VALUES(phone);

INSERT INTO addresses (id, user_id, recipient, phone, province, city, district, detail, postal_code, is_default)
VALUES
    ('addr_demo', 'usr_demo', 'Demo User', '18800000000', 'Guangdong', 'Shenzhen', 'Nanshan', 'Technology Park Center', '518000', TRUE)
ON DUPLICATE KEY UPDATE
    recipient = VALUES(recipient),
    phone = VALUES(phone),
    province = VALUES(province),
    city = VALUES(city),
    district = VALUES(district),
    detail = VALUES(detail),
    postal_code = VALUES(postal_code),
    is_default = VALUES(is_default);

UPDATE users SET default_address_id = 'addr_demo' WHERE id = 'usr_demo' AND (default_address_id IS NULL OR default_address_id = 'addr_demo');
