package util

type Codes struct {
	// #成功
	Success uint // 200
	// #失败
	Fail uint // 300

	//# 认证错误
	AuthErr uint // 401
	//MSG_AUTH_ERROR = 'token认证失败, 请重新登录'

	//# 服务器内部错误，状态码500
	ServerErr uint // = 500
	//MSG_SERVER_ERROR = '网络操作失败，请稍后重试'

	//# 未发现接口
	NotFoundErr uint // 404
	//MSG_NOT_FOUND_ERROR = '服务器没有此接口'

	//# 未知错误
	UnknownErr uint // 405
	//MSG_UNKNOWN_ERROR = '未知错误'

	//# 参数错误
	ParamErr uint // 201
	//MSG_PARAMETER_ERROR = '参数错误'

	//# 拒绝访问
	RejectErr uint // 202
	//MSG_REJECT_ERROR = '拒绝访问'

	//# 拒绝访问
	MethodErr uint // 203
	//MSG_METHOD_ERROR = '请求方法错误'

	//# 缺少参数
	ParamLack uint // 204
	//MSG_PARAMETER_LACK = '缺少参数'

	//# 业务上的错误
	UserExistsErr uint // 205
	//MSG_BUSSINESS_ERROR = '用户已存在'

	//# 查询失败
	QueryErr uint //206
	//MSG_QUERY_ERROR = '查询失败'

	//# 查询为空
	EmptyErr uint // 207
	//MSG_EMPTY_ERROR = '查询为空'

	//# 整顿期间
	CleanUp uint // 208
	//MSG_CLEAN_UP = '整顿期间'

	//# 未绑定手机号
	PhoneUnbind uint // 209
	//MSG_PHONE_UNBIND = '未绑定手机号'

	//# 手机验证
	PhoneUncheck uint //210
	//MSG_PHONE_UNCHECK = '未验证手机号'

	//# 不支持邮箱登录
	EmailErr uint // 211
	//MSG_EMAIL_ERROR = '暂不支持邮箱登录'

	PhoneUsed uint // 212
	//MSG_PHONE_USED = '该手机号已被使用'

	CreateErr   uint
	UserNotFont uint
}

var ApiCode = &Codes{
	Success:       200,
	Fail:          400,
	AuthErr:       401,
	ServerErr:     500,
	NotFoundErr:   404,
	UnknownErr:    420,
	ParamErr:      421,
	RejectErr:     422,
	MethodErr:     423,
	ParamLack:     424,
	UserExistsErr: 425,
	QueryErr:      426,
	EmptyErr:      427,
	CleanUp:       428,
	PhoneUnbind:   429,
	PhoneUncheck:  430,
	EmailErr:      431,
	PhoneUsed:     432,
	CreateErr:     433,
	UserNotFont:   434,
}

type Messages struct {
	// #成功
	Success string // 200

	// #失败
	Fail string // 300

	//# 认证错误
	AuthErr string // 401
	//MSG_AUTH_ERROR = 'token认证失败, 请重新登录'

	//# 服务器内部错误，状态码500
	ServerErr string // = 500
	//MSG_SERVER_ERROR = '网络操作失败，请稍后重试'

	//# 未发现接口
	NotFoundErr string // 404
	//MSG_NOT_FOUND_ERROR = '服务器没有此接口'

	//# 未知错误
	UnknownErr string // 405
	//MSG_UNKNOWN_ERROR = '未知错误'

	//# 参数错误
	ParamErr string // 201
	//MSG_PARAMETER_ERROR = '参数错误'

	//# 拒绝访问
	RejectErr string // 202
	//MSG_REJECT_ERROR = '拒绝访问'

	//# 拒绝访问
	MethodErr string // 203
	//MSG_METHOD_ERROR = '请求方法错误'

	//# 缺少参数
	ParamLack string // 204
	//MSG_PARAMETER_LACK = '缺少参数'

	//# 业务上的错误
	UserExistsErr string // 205
	//MSG_BUSSINESS_ERROR = '用户已存在'

	//# 查询失败
	QueryErr string //206
	//MSG_QUERY_ERROR = '查询失败'

	//# 查询为空
	EmptyErr string // 207
	//MSG_EMPTY_ERROR = '查询为空'

	//# 整顿期间
	CleanUp string // 208
	//MSG_CLEAN_UP = '整顿期间'

	//# 未绑定手机号
	PhoneUnbind string // 209
	//MSG_PHONE_UNBIND = '未绑定手机号'

	//# 手机验证
	PhoneUncheck string //210
	//MSG_PHONE_UNCHECK = '未验证手机号'

	//# 不支持邮箱登录
	EmailErr string // 211
	//MSG_EMAIL_ERROR = '暂不支持邮箱登录'

	PhoneUsed string // 212
	//MSG_PHONE_USED = '该手机号已被使用'
	CreateErr string
	// MSG_CREATE_ERR = '创建失败'
	UserNotFound string
}

var AMsg = &Messages{
	Success:       "成功",
	Fail:          "失败",
	AuthErr:       "token认证失败, 请重新登录",
	ServerErr:     "网络操作失败，请稍后重试",
	NotFoundErr:   "服务器没有此接口",
	UnknownErr:    "未知错误",
	ParamErr:      "参数错误",
	RejectErr:     "拒绝访问",
	MethodErr:     "请求方法错误",
	ParamLack:     "缺少参数",
	UserExistsErr: "UserExists",
	QueryErr:      "查询失败",
	EmptyErr:      "查询为空",
	CleanUp:       "整顿期间",
	PhoneUnbind:   "未绑定手机号",
	PhoneUncheck:  "未验证手机号",
	EmailErr:      "暂不支持邮箱登录",
	PhoneUsed:     "该手机号已被使用",
	CreateErr:     "创建失败",
	UserNotFound:  "用户不存在",
}
