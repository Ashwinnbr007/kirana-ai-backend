-- +goose Up
-- +goose StatementBegin
CREATE TABLE sales (
    item TEXT,
    quantity INTEGER,
    unit TEXT,
    retail_price_per_quantity DECIMAL,
    total_selling_price DECIMAL,
    CHECK (unit IN ('kg', 'g', 'dozen', 'unit'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sales;
-- +goose StatementEnd
