package role

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/model"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/do"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (s *Service) CreateRole(ctx context.Context, adminUser *common.AdminUser, req *dto.AddRoleReq) (int64, common.Errno) {
	roleID, err := s.adminRole.CreateRole(ctx, &do.AddRoleReq{
		AdminUserID: adminUser.ID,
		Name:        req.Name,
		Desc:        req.Desc,
	})
	if err != nil {
		logger.Error("CreateRole CreateRole error", zap.Error(err))
		return 0, common.DatabaseErr.WithErr(err)
	}
	return roleID, common.OK
}

func (s *Service) UpdateRole(ctx context.Context, adminUser *common.AdminUser, req *dto.UpdateRoleReq) common.Errno {
	err := s.adminRole.UpdateRole(ctx, &do.UpdateRoleReq{
		AdminUserID: adminUser.ID,
		ID:          req.ID,
		Name:        req.Name,
		Desc:        req.Desc,
		Status:      req.Status,
	})
	if err != nil {
		logger.Error("UpdateRole UpdateRole error", zap.Error(err))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) RoleList(ctx context.Context, req *dto.ListRoleReq) (*dto.ListRoleResp, common.Errno) {
	roleList, total, err := s.adminRole.ListRoles(ctx, &do.ListRoleReq{
		NameKw: req.NameKw,
		Status: req.Status,
		Pager:  req.Pager,
	})
	if err != nil {
		logger.Error("RoleList ListRoles error", zap.Error(err), zap.Any("req:", req))
		return nil, common.DatabaseErr.WithErr(err)
	}

	//获取角色id
	roleIds := make([]int64, 0)
	lo.ForEach(roleList, func(role *model.Role, index int) {
		roleIds = append(roleIds, role.ID)
	})
	//获取角色对应的权限
	rolePermMap, err := s.adminRole.GetRolePerms(ctx, lo.Uniq(roleIds))
	if err != nil {
		logger.Error("RoleList GetRolePerms error", zap.Error(err), zap.Any("role_ids:", roleIds))
		return nil, common.DatabaseErr.WithErr(err)
	}

	//获取权限id
	permIds := make([]int64, 0)
	for _, vList := range rolePermMap {
		permIds = append(permIds, vList...)
	}
	//获取权限名称
	permNameMap, err := s.adminPerm.GetPermNameMap(ctx, lo.Uniq(permIds))
	if err != nil {
		logger.Error("RoleList GetPermsByIds error", zap.Error(err), zap.Any("perm_ids:", permIds))
		return nil, common.DatabaseErr.WithErr(err)
	}

	retList := make([]*dto.RoleDto, 0, len(roleList))
	lo.ForEach(roleList, func(role *model.Role, index int) {
		perms := make([]common.IDName, 0)
		//rolePermMap[role.ID]是获取 权限id列表，然后遍历权限id列表
		lo.ForEach(rolePermMap[role.ID], func(permId int64, index int) {
			perms = append(perms, common.IDName{
				ID:   permId,
				Name: permNameMap[permId],
			})
		})
		retList = append(retList, &dto.RoleDto{
			ID:       role.ID,
			Name:     role.Name,
			Desc:     role.Desc,
			Status:   role.Status,
			Perms:    perms, // TODO 获取角色权限
			CreateAt: role.CreateAt.UnixMilli(),
			UpdateAt: role.UpdateAt.UnixMilli(),
		})
	})
	return &dto.ListRoleResp{
		Pager: req.Pager,
		List:  retList,
		Total: total,
	}, common.OK
}

func (s *Service) GetMyRoles(ctx context.Context, adminUser *common.AdminUser) ([]*dto.RoleDto, common.Errno) {
	//获取 用户对应的角色数据 admin_user_role
	list, err := s.adminRole.GetRolesByUserID(ctx, adminUser.ID)
	if err != nil {
		logger.Error("GetMyRoles GetMyRoles error", zap.Error(err), zap.Any("user_id:", adminUser.ID))
		return nil, common.DatabaseErr.WithErr(err)
	}

	//获取角色id
	roleIds := make([]int64, 0)
	lo.ForEach(list, func(role *model.AdminUserRole, index int) {
		roleIds = append(roleIds, role.RoleID)
	})
	//获取角色对应的权限
	rolePermMap, err := s.adminRole.GetRolePerms(ctx, lo.Uniq(roleIds))
	if err != nil {
		logger.Error("GetMyRoles GetRolePerms error", zap.Error(err), zap.Any("role_ids:", roleIds))
		return nil, common.DatabaseErr.WithErr(err)
	}

	//获取角色数据 roles
	roleMap, err := s.adminRole.GetRoleByIds(ctx, lo.Uniq(roleIds))
	if err != nil {
		logger.Error("GetMyRoles GetRoleByIds error", zap.Error(err), zap.Any("role_ids:", roleIds))
		return nil, common.DatabaseErr.WithErr(err)
	}

	//获取权限id
	permIds := make([]int64, 0)
	for _, vList := range rolePermMap {
		permIds = append(permIds, vList...)
	}
	//获取权限名称
	permNameMap, err := s.adminPerm.GetPermNameMap(ctx, lo.Uniq(permIds))
	if err != nil {
		logger.Error("GetMyRoles GetPermsByIds error", zap.Error(err), zap.Any("perm_ids:", permIds))
		return nil, common.DatabaseErr.WithErr(err)
	}

	retList := make([]*dto.RoleDto, 0, len(list))
	lo.ForEach(list, func(item *model.AdminUserRole, index int) {
		role, ok := roleMap[item.RoleID]
		if !ok {
			role = &model.Role{}
		}
		perms := make([]common.IDName, 0)
		//rolePermMap[item.ID]是获取 权限id列表，然后遍历权限id列表
		lo.ForEach(rolePermMap[item.ID], func(permId int64, index int) {
			perms = append(perms, common.IDName{
				ID:   permId,
				Name: permNameMap[permId],
			})
		})
		retList = append(retList, &dto.RoleDto{
			ID:       item.ID,
			Name:     role.Name,
			Desc:     role.Desc,
			Status:   role.Status,
			Perms:    perms,
			CreateAt: role.CreateAt.UnixMilli(),
			UpdateAt: role.UpdateAt.UnixMilli(),
		})
	})
	return retList, common.OK
}

func (s *Service) SetRolePerms(ctx context.Context, adminUser *common.AdminUser, req *dto.SetRolePermReq) common.Errno {
	err := s.adminRole.SetRolePerms(ctx, req.RoleID, req.PermIDs, adminUser.ID)
	if err != nil {
		logger.Error("SetRolePerms SetRolePerms error", zap.Error(err), zap.Any("req:", req))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}
