-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS role_permissions(
    ID uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    role_id uuid NOT NULL,
    permission_id uuid NOT NULL,
    CONSTRAINT fk_role_permission_permission_id FOREIGN KEY(permission_id) REFERENCES permissions(id),
    CONSTRAINT fk_role_permission_role_id FOREIGN KEY(role_id) REFERENCES roles(id),
    CONSTRAINT uc_role_permission_permission_id_role_id UNIQUE (role_id,permission_id)
);

INSERT INTO public.role_permissions (id,role_id,permission_id) VALUES
	('933d31f3-b100-47a2-b9d5-bc6589ce8aad','174ec139-715a-4b84-875a-53e081d6878c','00cd65c2-98c8-4475-bf3f-3c583469d919'),
	('693bc635-5ca1-4c57-993c-58007c8678f6','174ec139-715a-4b84-875a-53e081d6878c','b2047a0f-59c6-46d3-9eca-ba716abc656f'),
	('a238d7af-dbc0-40cc-8384-14a70f210d73','174ec139-715a-4b84-875a-53e081d6878c','bf4e63f8-9273-43d2-a4b6-c7cfa4861cf7'),
	('18a571b6-1ae1-4fff-8ae7-94ba494b3b77','174ec139-715a-4b84-875a-53e081d6878c','e0771b2c-f02f-4c6b-a050-5c1c0959e92d'),
	('246cf5b4-b3fb-45be-830e-057d0f92a4fd','c98c4ef9-f6f6-48af-a9f7-98f57aae0143','44be6cd5-544b-4b89-90cf-1973ee4c7a5f'),
	('23157781-2d75-47a7-a9f0-281a5d321ae7','c98c4ef9-f6f6-48af-a9f7-98f57aae0143','b2047a0f-59c6-46d3-9eca-ba716abc656f');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS role_permissions;
-- +goose StatementEnd
