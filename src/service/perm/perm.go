package perm

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/do"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (s *Service) CreatePermission(ctx context.Context, user *common.AdminUser, req *dto.AddPermissionReq) (int64, common.Errno) {
	permID, err := s.adminPerm.CreatePermission(ctx, &do.AddPerm{
		AdminUserID: user.ID,
		Code:        req.Code,
		Name:        req.Name,
		PagePath:    req.PagePath,
		ParentID:    req.ParentID,
		Sort:        req.Sort,
		Type:        req.Type,
		Desc:        req.Desc,
	})
	if err != nil {
		logger.Error("CreatePermission CreatePermission error", zap.Error(err), zap.Any("req", req))
		return 0, common.DatabaseErr.WithErr(err)
	}
	return permID, common.OK
}

func (s *Service) UpdatePermissions(ctx context.Context, adminUser *common.AdminUser, req *dto.UpdatePermissionReq) common.Errno {
	reqList := make([]do.UpdatePerm, 0)
	lo.ForEach(req.List, func(item dto.UpdatePermDto, index int) {
		reqList = append(reqList, do.UpdatePerm{
			ID: item.ID,
			AddPerm: do.AddPerm{
				AdminUserID: adminUser.ID,
				Code:        item.Code,
				Name:        item.Name,
				PagePath:    item.PagePath,
				ParentID:    lo.Ternary(item.ParentID == 0, -1, item.ParentID),
				Sort:        item.Sort,
				Type:        item.Type,
				Desc:        item.Desc,
			},
		})
	})
	err := s.adminPerm.UpdatePermission(ctx, &do.UpdatePermList{
		List: reqList,
	})
	if err != nil {
		logger.Error("UpdatePermissions UpdatePermission error", zap.Error(err), zap.Any("req", req))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) DeletePermission(ctx context.Context, adminUser *common.AdminUser, req *dto.DeletePermissionReq) common.Errno {
	err := s.adminPerm.DeletePermission(ctx, &do.DeletePerm{
		ID: req.ID,
	})
	if err != nil {
		logger.Error("DeletePermission DeletePermission error", zap.Error(err), zap.Any("req", req))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}
