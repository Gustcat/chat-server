-- +goose Up
CREATE TABLE chat_users(
    chat_id int not null,
    username varchar(50) not null,
    created_at timestamp not null default now(),
    CONSTRAINT pk_chat_username PRIMARY KEY (chat_id, username),
    CONSTRAINT fk_chat_users_chat FOREIGN KEY (chat_id) REFERENCES chat(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chat_users;