-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS counters (
    name TEXT PRIMARY KEY,
    value BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS gauges (
    name TEXT PRIMARY KEY,
    value DOUBLE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gauges;

DROP TABLE counters;
-- +goose StatementEnd
