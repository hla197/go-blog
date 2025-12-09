package errors

const (
	INVALID_PARAMETER int = 1001 + iota // 参数错误
	SYSTEM_ERROR                        //系统错误
	OTHER_ERROR
)

const (
	AUTH_ERROR int = 2001 + iota
	POST_ERROR
	COMMENT_ERROR
)
