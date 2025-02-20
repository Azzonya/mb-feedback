do
$$
    begin
        execute 'ALTER DATABASE ' || current_database() || ' SET timezone = ''+05''';
    end;
$$;

CREATE TABLE IF NOT EXISTS ord (
    id BIGSERIAL PRIMARY KEY,
    external_order_id VARCHAR(50) NOT NULL UNIQUE,  -- ID заказа во внешнем сервисе
    user_phone VARCHAR(20) NOT NULL,               -- номер телефона покупателя
    user_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

 CREATE TABLE IF NOT EXISTS ord_detail (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES ord (id),
    product_code VARCHAR(100) NOT NULL,  -- код товара
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

 CREATE TABLE IF NOT EXISTS notification (
    id BIGSERIAL PRIMARY KEY,
    order_item_id BIGINT NOT NULL REFERENCES ord_detail(id),
    phone_number VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    sent_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

