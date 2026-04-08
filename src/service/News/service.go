package News

import (
	"crypto/rand"
	"math/big"
	"strconv"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/redis"
	"github.com/Sam-Stranding/SamMall/src/config"
)

type SmsCode struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	ErrCode int64  `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

type Service struct {
	conf   *config.Config
	verify redis.IVerify
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		conf:   adaptor.GetConf(),
		verify: redis.NewVerify(adaptor),
	}
}

// 随机生成验证码
func GenerateRandomNumber() string {
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(n.Int64()) + 1000)
}
