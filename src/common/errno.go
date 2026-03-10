package common

type Errno struct {
	Code   int
	Msg    string
	ErrMsg string //调试时使用
}

func (err Errno) Error() string {
	return err.Msg
}

func (err Errno) WithMsg(msg string) string {
	err.Msg = err.Msg + "," + msg
	return err.Msg
}

func (err Errno) WithErr(rawErr error) Errno {
	var msg string
	if rawErr != nil {
		msg = rawErr.Error()
	}
	err.ErrMsg = err.Msg + "," + msg
	return err
}

var (
	OK            = Errno{Code: 200, Msg: "OK"}
	ServerErr     = Errno{Code: 500, Msg: "Internal Server Error"}
	ParamErr      = Errno{Code: 400, Msg: "Param Error"}
	AuthErr       = Errno{Code: 401, Msg: "Auth Error"}
	PermissionErr = Errno{Code: 403, Msg: "Permission Denied"}

	DatabaseErr = Errno{Code: 10000, Msg: "Database Error"}
	RedisErr    = Errno{Code: 10001, Msg: "Redis Error"}

	UserNotFoundErr = Errno{Code: 11001, Msg: "User Exist"}
)
