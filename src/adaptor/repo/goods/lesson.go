package goods

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/model"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/query"
	"github.com/Sam-Stranding/SamMall/src/service/do"
	"gorm.io/gorm"
)

type ILesson interface {
	AddCategory(ctx context.Context, req *do.AddCategory) (int64, error)
	UpdateCategory(ctx context.Context, req *do.UpdateCategory) error
	DeleteCategory(ctx context.Context, req *do.DeleteCategory) error
	ListCategory(ctx context.Context, req *do.ListCategory) ([]*model.LessonCategory, error)
	UpdateSort(ctx context.Context, sortList []*do.UpdateSort) error
}

type Lesson struct {
	db *gorm.DB
}

func NewLesson(adaptor adaptor.IAdaptor) *Lesson {
	return &Lesson{
		db: adaptor.GetDB(),
	}
}

func (l *Lesson) AddCategory(ctx context.Context, req *do.AddCategory) (int64, error) {
	qs := query.Use(l.db).LessonCategory
	addObject := &model.LessonCategory{
		Name:     req.Name,
		Level:    req.Level,
		ParentID: req.ParentID,
		Sort:     req.Sort,
	}
	err := qs.WithContext(ctx).Create(addObject)
	return addObject.ID, err
}

func (l *Lesson) UpdateCategory(ctx context.Context, req *do.UpdateCategory) error {
	qs := query.Use(l.db).LessonCategory
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Update(qs.Name, req.Name)
	return err
}

func (l *Lesson) DeleteCategory(ctx context.Context, req *do.DeleteCategory) error {
	qs := query.Use(l.db).LessonCategory
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Delete()
	return err
}

func (l *Lesson) ListCategory(ctx context.Context, req *do.ListCategory) ([]*model.LessonCategory, error) {
	qs := query.Use(l.db).LessonCategory
	list, err := qs.WithContext(ctx).Order(qs.Sort).Find()
	return list, err
}

func (l *Lesson) UpdateSort(ctx context.Context, sortList []*do.UpdateSort) error {
	qs := query.Use(l.db).LessonCategory
	return l.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, v := range sortList {
			_, err := qs.WithContext(ctx).Where(qs.ID.Eq(v.ID)).Update(qs.Sort, v.Sort)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
