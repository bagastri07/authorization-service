-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS permissions (
    ID uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL
);

INSERT INTO public.permissions (id,"name") VALUES
    ('00cd65c2-98c8-4475-bf3f-3c583469d919','CREATE_PRODUCT'),
    ('e0771b2c-f02f-4c6b-a050-5c1c0959e92d','UPDATE_PRODUCT'),
    ('bf4e63f8-9273-43d2-a4b6-c7cfa4861cf7','DELETE_PRODUCT'),
    ('b2047a0f-59c6-46d3-9eca-ba716abc656f','READ_PRODUCT'),
    ('44be6cd5-544b-4b89-90cf-1973ee4c7a5f','BUY_PRODUCT');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS permissions;
-- +goose StatementEnd
