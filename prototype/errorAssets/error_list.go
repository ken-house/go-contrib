package errorAssets

var (
	// 在项目目录下（internal/lib/errorAssets/）创建error_list.go文件；
	// 通过下面示例代码定义错误码
	ERR_PARAM     = NewError(10000, "参数错误")
	ERR_SYSTEM    = NewError(10001, "系统错误")
	ERR_CERT      = NewError(10002, "证书错误")
	ERR_CALL_FUNC = NewError(10003, "调用方法出错")
	ERR_DIAL      = NewError(10004, "连接错误")
)
