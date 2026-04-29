package admin

import (
	"context"
	"testing"
	"time"

	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSetRolePermsReplacesExistingPermissions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	err = db.AutoMigrate(&model.RolePermission{})
	if err != nil {
		t.Fatalf("migrate role_permission: %v", err)
	}

	now := time.Now()
	seed := []*model.RolePermission{
		{
			RoleID:       7,
			PermissionID: 10,
			CreateAt:     now,
			UpdateAt:     now,
			CreateBy:     1,
			UpdateBy:     1,
		},
		{
			RoleID:       7,
			PermissionID: 12,
			CreateAt:     now,
			UpdateAt:     now,
			CreateBy:     1,
			UpdateBy:     1,
		},
	}
	if err := db.Create(seed).Error; err != nil {
		t.Fatalf("seed role permissions: %v", err)
	}

	repo := &AdminRole{db: db}
	if err := repo.SetRolePerms(context.Background(), 7, []int64{21}, 2); err != nil {
		t.Fatalf("set role perms: %v", err)
	}

	permMap, err := repo.GetRolePerms(context.Background(), []int64{7})
	if err != nil {
		t.Fatalf("get role perms: %v", err)
	}

	got := permMap[7]
	if len(got) != 1 {
		t.Fatalf("expected 1 permission after replacement, got %v", got)
	}
	if got[0] != 21 {
		t.Fatalf("expected permission 21 after replacement, got %v", got)
	}
}
