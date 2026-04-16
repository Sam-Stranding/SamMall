package admin

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/model"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/do"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *Service) GetAdminUserByToken(ctx context.Context, token string) (*common.AdminUser, common.Errno) {
	userString, err := s.verify.GetAdminUserToken(ctx, token)
	if err != nil {
		logger.Error("GetAdminUserToken err", zap.Error(err), zap.Any("token", token))
		return nil, common.RedisErr.WithErr(err)
	}
	adminUser := &common.AdminUser{}
	err = json.Unmarshal([]byte(userString), adminUser)
	if err != nil {
		logger.Error("GetAdminUserToken json.Unmarshal err", zap.Error(err), zap.String("userString", userString))
		return nil, common.RedisErr.WithErr(err)
	}
	return adminUser, common.OK
}

func (s *Service) AdminUserList(ctx context.Context, adminUser *common.AdminUser, req *dto.ListAdminUserReq) (*dto.ListAdminUserResp, common.Errno) {
	userList, total, err := s.adminUser.ListAdminUser(ctx, &do.ListAdminUser{
		Name:   req.Name,
		Mobile: req.Mobile,
		RoleID: req.RoleID,
		Status: req.Status,
		Pager:  req.Pager,
	})
	if err != nil {
		logger.Error("AdminUserList ListAdminUser err", zap.Error(err), zap.Any("req", req))
		return nil, common.DatabaseErr.WithErr(err)
	}

	userIds := make([]int64, 0)
	lo.ForEach(userList, func(item *model.AdminUser, index int) {
		userIds = append(userIds, item.ID)
	})
	//获取用户对应的role_id
	userRoleMap, err := s.adminRole.GetRoleByUserIds(ctx, userIds)
	if err != nil {
		logger.Error("AdminUserList GetRoleByUserIds err", zap.Error(err), zap.Any("user_ids", userIds))
		return nil, common.DatabaseErr.WithErr(err)
	}

	roleIds := make([]int64, 0)
	for _, vList := range userRoleMap {
		for _, v := range vList {
			roleIds = append(roleIds, v.ID)
		}
	}
	//通过role_id获取name(权限名称)
	roleMap, err := s.adminRole.GetRoleByIds(ctx, lo.Uniq(roleIds))
	if err != nil {
		logger.Error("AdminUserList GetRoleByIds err", zap.Error(err), zap.Any("role_ids", roleIds))
		return nil, common.DatabaseErr.WithErr(err)
	}

	retList := make([]*dto.AdminUserWithRoleDto, 0)
	lo.ForEach(userList, func(user *model.AdminUser, index int) {
		retList = append(retList, &dto.AdminUserWithRoleDto{
			AdminUserDto: dto.AdminUserDto{
				ID:         user.ID,
				UserID:     user.ID,
				Name:       user.Name,
				NickName:   user.NickName,
				Mobile:     user.Mobile,
				Sex:        user.Sex,
				Status:     user.Status,
				LarkOpenID: user.LarkOpenID,
				CreateAt:   user.CreateAt.UnixMilli(),
				UpdateAt:   user.UpdateAt.UnixMilli(),
			},
			Roles: lo.Map(userRoleMap[user.ID], func(item *model.AdminUserRole, index int) *common.IDName {
				return &common.IDName{
					ID:   item.RoleID,
					Name: roleMap[item.RoleID].Name,
				}
			}),
		})
	})
	return &dto.ListAdminUserResp{
		List:  retList,
		Total: total,
	}, common.OK
}

func (s *Service) GetUserInfo(ctx context.Context, adminUser *common.AdminUser) (*dto.AdminUserWithRoleDto, common.Errno) {
	user, err := s.adminUser.GetUserInfo(ctx, adminUser.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.UserNotFoundErr
		}
		logger.Error("GetUserInfo GetUserInfo err", zap.Error(err), zap.Any("user_id", adminUser))
		return nil, common.DatabaseErr.WithErr(err)
	}
	userIds := make([]int64, 0)
	userIds = append(userIds, user.ID)
	//获取用户对应的role_id
	userRoleMap, err := s.adminRole.GetRoleByUserIds(ctx, userIds)
	if err != nil {
		logger.Error("AdminUserList GetRoleByUserIds err", zap.Error(err), zap.Any("user_ids", userIds))
		return nil, common.DatabaseErr.WithErr(err)
	}

	roleIds := make([]int64, 0)
	for _, vList := range userRoleMap {
		for _, v := range vList {
			roleIds = append(roleIds, v.ID)
		}
	}
	//通过role_id获取name(权限名称)
	roleMap, err := s.adminRole.GetRoleByIds(ctx, lo.Uniq(roleIds))
	if err != nil {
		logger.Error("AdminUserList GetRoleByIds err", zap.Error(err), zap.Any("role_ids", roleIds))
		return nil, common.DatabaseErr.WithErr(err)
	}

	ret := &dto.AdminUserWithRoleDto{
		AdminUserDto: dto.AdminUserDto{
			ID:         user.ID,
			UserID:     user.ID,
			Name:       user.Name,
			NickName:   user.NickName,
			Mobile:     user.Mobile,
			Sex:        user.Sex,
			Status:     user.Status,
			LarkOpenID: user.LarkOpenID,
			CreateAt:   user.CreateAt.UnixMilli(),
		},
		Roles: lo.Map(userRoleMap[user.ID], func(item *model.AdminUserRole, index int) *common.IDName {
			return &common.IDName{
				ID:   item.RoleID,
				Name: roleMap[item.RoleID].Name,
			}
		}),
	}

	return ret, common.OK
}

func (s *Service) CreateUser(ctx context.Context, adminUser *common.AdminUser, req *dto.CreateUserReq) (int64, common.Errno) {
	userID, err := s.adminUser.CreateUser(ctx, &do.CreateUser{
		Name:        req.Name,
		NickName:    req.NickName,
		Mobile:      req.Mobile,
		Sex:         req.Sex,
		AdminUserID: adminUser.UserID,
		RoleIDs:     req.RoleIDs,
	})
	if err != nil {
		logger.Error("Create Error", zap.Error(err), zap.Any("req", req))
		return 0, common.DatabaseErr.WithErr(err)
	}
	return userID, common.OK
}

func (s *Service) UpdateUser(ctx context.Context, adminUser *common.AdminUser, req *dto.UpdateUserReq) common.Errno {
	err := s.adminUser.UpdateUser(ctx, &do.UpdateUser{
		ID:          req.ID,
		Name:        req.Name,
		NickName:    req.NickName,
		Sex:         req.Sex,
		AdminUserID: adminUser.UserID,
		Status:      req.Status,
		RoleIDs:     req.RoleIDs,
	})
	if err != nil {
		logger.Error("UpdateUser Error", zap.Error(err), zap.Any("req", req))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) LarkBind(ctx context.Context, adminUser *common.AdminUser, req *dto.LarkQrCodeBindReq) common.Errno {
	accessToken, errno := s.token.GetLarkUserAccessToken(ctx, req.AppCode, req.Code, req.RedirectUri, "", false)
	if errno.NotOk() {
		logger.Error("LarkBind GetLarkUserAccessToken Error", zap.Error(errno), zap.Any("req", req))
		return common.ServerErr.WithErr(errno)
	}
	larkUserInfo, err := s.lark.GetLarkUserInfo(ctx, accessToken.Token)
	if err != nil {
		logger.Error("LarkBind GetLarkUserInfo Error", zap.Error(err), zap.Any("req", req))
		return common.ServerErr.WithErr(err)
	}
	err = s.adminUser.UpdateUserLarkOpenID(ctx, adminUser.UserID, larkUserInfo.OpenID)
	if err != nil {
		logger.Error("LarkBind UpdateUserLarkOpenID Error", zap.Error(err), zap.Any("adminUser", adminUser))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) LarkUnbind(ctx context.Context, adminUser *common.AdminUser) common.Errno {
	err := s.adminUser.UpdateUserLarkOpenID(ctx, adminUser.UserID, "")
	if err != nil {
		logger.Error("LarkUnbind UpdateUserLarkOpenID Error", zap.Error(err), zap.Any("adminUser", adminUser))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) AdminUserLogout(ctx context.Context, adminUser *common.AdminUser) common.Errno {
	err := s.verify.CleanToken(ctx, adminUser.UserID)
	if err != nil {
		logger.Error("AdminUserLogout CleanToken Error", zap.Error(err), zap.Any("adminUser", adminUser))
		return common.RedisErr.WithErr(err)
	}
	return common.OK
}
