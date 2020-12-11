
-- +migrate Up
ALTER TABLE
    `users`
ADD
    `status` VARCHAR(32) NOT NULL AFTER `created_at`,
ADD
    `password` VARCHAR(100) NOT NULL AFTER `status`,
CHANGE
    `created_at` `created_at` DATETIME  NOT NULL;

-- +migrate Down
ALTER TABLE
    `users`
DROP
    `status`,
DROP
    `password`,
CHANGE
    `created_at` `created_at` varchar(100) NOT NULL;
