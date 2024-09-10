package errno

// 基础错误码
var (
	OK                  = &Errno{Code: 0, Msg: "OK"}
	BadRequest          = &Errno{Code: 400, Msg: "Bad request"}
	Unauthorized        = &Errno{Code: 401, Msg: "Unauthorized"}
	Forbidden           = &Errno{Code: 403, Msg: "Forbidden"}
	InternalServerError = &Errno{Code: 500, Msg: "Internal server error"}
)

type Errno struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (err *Errno) Error() string {
	return err.Msg
}

func (err *Errno) SetMessage(msg string) *Errno {
	err.Msg = msg
	return err
}

// EncodeErr decode error to Errno
func EncodeErr(err error) *Errno {
	if err == nil {
		return OK
	}
	val, ok := err.(*Errno)
	if !ok {
		return InternalServerError
	}
	return val
}
