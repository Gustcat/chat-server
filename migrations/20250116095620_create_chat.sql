-- +goose Up
CREATE TABLE chat (
      id serial primary key,
      created_at timestamp not null default now(),
      updated_at timestamp
);

-- +goose Down
DROP TABLE chat;
