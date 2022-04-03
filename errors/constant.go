package errors

const (
	SuccessCode = 200 // 成功

	// UnauthorizedCode 服务请求相关
	UnauthorizedCode   = 40401 // 未授权
	SignatureCode      = 40402 // 签名不匹配
	CryptCode          = 40403 // 加解密错误
	NotFoundCode       = 40404 // 未找到
	TooManyRequestCode = 40407 // 请求太频繁
	RequestTimeoutCode = 40408 // 请求超时

	// ServerUnavailableCode 服务响应相关
	ServerUnavailableCode = 40500 // rpc服务异常
	DatabaseCode          = 40510 // 数据库异常
	InvalidParamCode      = 40520 // 请求参数错误

	OtherErrorCode = 50001

	MysqlErrUnsignedOutOfRange = 1690 // mysql Error 1690: BIGINT UNSIGNED value is out of range in

	ErrTokenInvalid = 40401 // Token无效

	ErrTimeout        = 40501 // 服务超时
	ErrServiceInvalid = 40503 // 服务调用错误

	ErrInternal      = 40500 // 服务器内部错误
	ErrParamsInvalid = 40520 // 请求参数错误

	ErrMarshalError   = 40530 // 参数Marshal失败
	ErrUnMarshalError = 40531 // 参数UnMarshal失败

)

const (
	ErrMsgSuccess           = "成功"
	ErrMsgInvalidParam      = "请求参数错误"
	ErrMsgNotFount          = "数据未找到"
	ErrMsgServerUnavailable = "服务开小差了，请稍后再试"
	ErrMsgTooManyRequest    = "请求太频繁，请稍后再试"
	ErrMsgRequestTimeout    = "请求超时，请稍后再试"
	ErrMsgCrypt             = "数据加解密异常"
	ErrMsgSignature         = "签名不匹配"

	ErrMsgUnauthorized = "授权失败，请登录"
)

var (
	ErrInvalidToken      = New(UnauthorizedCode, ErrMsgUnauthorized)
	ErrNotFound          = New(NotFoundCode, ErrMsgNotFount)
	ErrTooManyRequest    = New(TooManyRequestCode, ErrMsgTooManyRequest)
	ErrRequestTimeout    = New(RequestTimeoutCode, ErrMsgRequestTimeout)
	ErrServerUnavailable = New(ServerUnavailableCode, ErrMsgServerUnavailable)
	ErrDatabase          = New(DatabaseCode, ErrMsgServerUnavailable)
	ErrInvalidParams     = New(InvalidParamCode, ErrMsgInvalidParam)
)

const (
	ServerErr = iota + 1
	JSONErr
	ParamsErr
)

type Err struct {
	Code    int32
	Message string
}

var Errs = map[int]*Err{
	ServerUnavailableCode: {501, ErrMsgServerUnavailable},
	JSONErr:               {406, "Json 格式错误"},
	ParamsErr:             {406, "参数错误"},
}
