package goods

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

func (s *Service) CreateCategory(ctx context.Context, req *dto.AddCategoryReq) (int64, common.Errno) {
	cate := &do.AddCategory{
		Name:     req.Name,
		Level:    req.Level,
		ParentID: req.ParentID,
		Sort:     req.Sort,
	}
	cateID, err := s.lesson.AddCategory(ctx, cate)
	if err != nil {
		logger.Error("CreateCategory AddCategory error", zap.Error(err), zap.Any("cate", cate))
		return 0, common.DatabaseErr.WithErr(err)
	}
	return cateID, common.OK
}

func (s *Service) UpdateCategory(ctx context.Context, req *dto.UpdateCategoryReq) common.Errno {
	cate := &do.UpdateCategory{
		ID:   req.ID,
		Name: req.Name,
	}
	err := s.lesson.UpdateCategory(ctx, cate)
	if err != nil {
		logger.Error("UpdateCategory UpdateCategory error", zap.Error(err), zap.Any("cate", cate))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) DeleteCategory(ctx context.Context, req *dto.DeleteCategoryReq) common.Errno {
	cate := &do.DeleteCategory{
		ID: req.ID,
	}
	err := s.lesson.DeleteCategory(ctx, cate)
	if err != nil {
		logger.Error("DeleteCategory DeleteCategory error", zap.Error(err), zap.Any("cate", cate))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) CategoryList(ctx context.Context, req *dto.ListCategoryReq) ([]*dto.CategoryDto, common.Errno) {
	list, err := s.lesson.ListCategory(ctx, &do.ListCategory{
		Pager: req.Pager,
	})
	if err != nil {
		logger.Error("CategoryList ListCategory error", zap.Error(err))
		return nil, common.DatabaseErr.WithErr(err)
	}
	return lo.Map(list, func(item *model.LessonCategory, index int) *dto.CategoryDto {
		return &dto.CategoryDto{
			ID:       item.ID,
			Name:     item.Name,
			Level:    item.Level,
			ParentID: item.ParentID,
			Sort:     item.Sort,
		}
	}), common.OK

}

func (s *Service) CategorySorts(ctx context.Context, sortList []dto.UpdateSort) common.Errno {
	updateSorts := make([]*do.UpdateSort, 0)
	lo.ForEach(sortList, func(item dto.UpdateSort, index int) {
		updateSorts = append(updateSorts, &do.UpdateSort{
			ID:   item.ID,
			Sort: item.Sort,
		})
	})
	err := s.lesson.UpdateSort(ctx, updateSorts)
	if err != nil {
		logger.Error("CategorySorts UpdateSort error", zap.Error(err))
		return common.DatabaseErr.WithErr(err)
	}
	return common.OK
}
