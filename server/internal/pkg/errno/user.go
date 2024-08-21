package errno

var (
	UserAlreadyExist = &Errno{
		Code: 401001,
		Msg:  "用户已存在",
	}
)
