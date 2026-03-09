package admin

import "github.com/Sam-Stranding/SamMall/src/adaptor"

type Service struct {
	adaptor *adaptor.Adaptor
}

func NewService(adaptor *adaptor.Adaptor) *Service {
	return &Service{
		adaptor: adaptor,
	}
}
