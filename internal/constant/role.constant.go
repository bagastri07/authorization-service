package constant

import "github.com/google/uuid"

type RoleID uuid.UUID

var (
	CustomerRoleID RoleID = RoleID(uuid.MustParse("c98c4ef9-f6f6-48af-a9f7-98f57aae0143"))
	AdminRoleID    RoleID = RoleID(uuid.MustParse("174ec139-715a-4b84-875a-53e081d6878c"))
)
