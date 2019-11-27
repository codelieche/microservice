package repositories

import "testing"

func TestReadyData(t *testing.T) {
	//	创建User、Group、Role、Application、Permission
	t.Run("create_users", TestUserRepository_Save)
	t.Run("create_groups", TestGroupRepository_Save)
	t.Run("create_roles", TestRoleRepository_Save)
	t.Run("create_apps", TestApplicationRepository_Save)
	t.Run("create_permissions", TestPermissionRepository_Save)
	t.Run("create_ticket", TestTicketRepository_Save)
}
