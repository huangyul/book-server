package errno

var (
	UserAlreadyExist = &Errno{
		Code: 401001,
		Msg:  "用户已存在",
	}
	UserNotFound = &Errno{
		Code: 401002,
		Msg:  "用户不存在",
	}
)
