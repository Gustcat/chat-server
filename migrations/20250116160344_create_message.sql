-- +goose Up
CREATE TABLE message (
    id serial primary key,
    text text,
    username varchar(50),
    chat_id int not null ,
    created_at timestamp not null default now(),
    CONSTRAINT fk_message_chat FOREIGN KEY (chat_id) REFERENCES chat(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE message;