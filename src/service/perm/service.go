package perm

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/admin"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/model"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Service struct {
	adminPerm admin.IPerm
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminPerm: admin.NewAdminPerm(adaptor),
	}
}

func (s *Service) PermissionList(ctx context.Context) (*dto.PermissionListResp, common.Errno) {
	permList, total, err := s.adminPerm.PermissionList(ctx, common.Pager{
		Page:      1,
		Limit:     1000,
		UnLimited: true,
	})
	if err != nil {
		logger.Error("PermissionList PermissionList error", zap.Error(err))
		return nil, common.DatabaseErr.WithErr(err)
	}
	retList := make([]*dto.PermissionDto, 0, len(permList))
	lo.ForEach(permList, func(item *model.Permission, index int) {
		retList = append(retList, &dto.PermissionDto{
			ID:       item.ID,
			Code:     item.Code,
			Name:     item.Name,
			Desc:     item.Desc,
			PagePath: item.PagePath,
			ParentID: item.ParentID,
			Sort:     item.Sort,
			Status:   item.Status,
			Type:     item.Type,
		})
	})
	return &dto.PermissionListResp{
		Pager: common.Pager{
			Page:      1,
			Limit:     1000,
			UnLimited: true,
		},
		Total: total,
		List:  retList,
	}, common.OK
}

func (s *Service) MyPermissionList(ctx context.Context, user *common.AdminUser) ([]*dto.PermissionDto, common.Errno) {
	permList, err := s.adminPerm.MyPermissionList(ctx, common.Pager{
		Page:      1,
		Limit:     1000,
		UnLimited: true,
	})
	if err != nil {
		logger.Error("MyPermissionList MyPermissionList error", zap.Error(err))
		return nil, common.DatabaseErr.WithErr(err)
	}
	retList := make([]*dto.PermissionDto, 0, len(permList))
	lo.ForEach(permList, func(item *model.Permission, index int) {
		retList = append(retList, &dto.PermissionDto{
			ID:       item.ID,
			Code:     item.Code,
			Name:     item.Name,
			Desc:     item.Desc,
			PagePath: item.PagePath,
			ParentID: item.ParentID,
			Sort:     item.Sort,
			Status:   item.Status,
			Type:     item.Type,
		})
	})
	return retList, common.OK
}
