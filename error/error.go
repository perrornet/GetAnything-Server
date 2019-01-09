package error

var (
	NotFoundHandlerError = &newError{msg: "没有寻找到能够处理当前url的程序，我们已经记录该url，将会尽快提供支持。"}
	ClientError          = &newError{msg: "客户端发送请求出现错误，我们将尽快解决！"}
	ServerError          = &newError{msg: "服务器出现错误，我们将尽快解决！"}
)

type newError struct{ msg string }

func (u *newError) Error() string { return u.msg }
