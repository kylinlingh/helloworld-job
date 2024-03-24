package errorcode

// 描述的是业务错误码（不能为 10 开头），注意注释一定要符合规范，才能保证自动生成代码：// <错误码整型常量> - <对应的HTTP Status Code>: .

// user errors.
const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 110001

	// ErrUserAlreadyExist - 400: User already exist.
	ErrUserAlreadyExist
)
