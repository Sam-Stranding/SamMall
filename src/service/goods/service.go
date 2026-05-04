package goods

import (
	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/goods"
)

type Service struct {
	lesson goods.ILesson
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		lesson: goods.NewLesson(adaptor),
	}
}
