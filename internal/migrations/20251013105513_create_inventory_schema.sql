-- +goose Up
-- +goose StatementBegin
CREATE TABLE inventory (
    item TEXT,
    quantity INTEGER,
    unit TEXT,
    wholesale_price_per_quantity INTEGER,
    total_cost_of_product INTEGER,
    CHECK (unit IN ('kg', 'g', 'dozen', 'unit'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE inventory;
-- +goose StatementEnd
