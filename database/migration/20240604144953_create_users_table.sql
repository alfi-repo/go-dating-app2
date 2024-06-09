-- +goose Up
CREATE TABLE `users`
(
    `id`           int PRIMARY KEY AUTO_INCREMENT,
    `email`        varchar(100) UNIQUE NOT NULL,
    `password`     varchar(255)        NOT NULL,
    `created_at`   timestamp           NOT NULL,
    `updated_at`   timestamp           NOT NULL,
    `suspended_at` timestamp
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE `users`;
-- +goose StatementBegin
-- +goose StatementEnd
