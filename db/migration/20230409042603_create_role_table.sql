-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    ID uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL
);

INSERT INTO public.roles (id,"name") VALUES
    ('174ec139-715a-4b84-875a-53e081d6878c','ADMIN'),
    ('c98c4ef9-f6f6-48af-a9f7-98f57aae0143','CUSTOMER');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
