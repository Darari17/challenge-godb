CREATE TABLE customer(
    customer_id INT PRIMARY KEY, 
    name VARCHAR(255) NOT NULL, 
    phone VARCHAR(255) NOT NULL, 
    address VARCHAR(255) DEFAULT '', 
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);