package errno

var (
	// UserAlreadyExist 用户已存在
	UserAlreadyExist = &Errno{
		Code: 401001,
		Msg:  "用户已存在",
	}
	// UserNotFound 用户不存在
	UserNotFound = &Errno{
		Code: 401002,
		Msg:  "用户不存在",
	}
	// UserPasswordIncorrect 账号或密码不正确
	UserPasswordIncorrect = &Errno{
		Code: 401003,
		Msg:  "账号或密码不正确",
	}
)
