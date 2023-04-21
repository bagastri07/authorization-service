-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_roles(
    ID uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    role_id uuid NOT NULL,
    user_id uuid NOT NULL,
    CONSTRAINT fk_user_role_role_id FOREIGN KEY(role_id) REFERENCES roles(id),
    CONSTRAINT fk_user_role_user_id FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT uc_user_role_user_id_role_id UNIQUE (role_id,user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_roles;
-- +goose StatementEnd
