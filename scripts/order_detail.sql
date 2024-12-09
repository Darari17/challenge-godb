CREATE TABLE order_detail (
    order_detail_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    service_id INT NOT NULL,
    qty INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES "order"(order_id),
    FOREIGN KEY (service_id) REFERENCES service(service_id)
);