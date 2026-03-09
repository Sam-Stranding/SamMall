package customer

import "github.com/Sam-Stranding/SamMall/src/adaptor"

type Ctrl struct {
	adaptor *adaptor.Adaptor
}

func NewCtrl(adaptor *adaptor.Adaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
	}
}
