package model

// ErrorCode 枚举领域内使用的错误标识。
type ErrorCode string

const (
	ErrCodeInvalidParameter ErrorCode = "ERR_INVALID_PARAMETER"
	ErrCodeUserNotFound     ErrorCode = "ERR_USER_NOT_FOUND"
	ErrCodeAddressNotFound  ErrorCode = "ERR_ADDRESS_NOT_FOUND"
	ErrCodeProductNotFound  ErrorCode = "ERR_PRODUCT_NOT_FOUND"
	ErrCodeInventoryShort   ErrorCode = "ERR_INVENTORY_SHORTAGE"
	ErrCodeCartEmpty        ErrorCode = "ERR_CART_EMPTY"
	ErrCodeOrderNotFound    ErrorCode = "ERR_ORDER_NOT_FOUND"
	ErrCodePaymentFailed    ErrorCode = "ERR_PAYMENT_FAILED"
)

// KnownErrorCodes 方便在文档接口中暴露支持的错误码。
var KnownErrorCodes = []ErrorCode{
	ErrCodeInvalidParameter,
	ErrCodeUserNotFound,
	ErrCodeAddressNotFound,
	ErrCodeProductNotFound,
	ErrCodeInventoryShort,
	ErrCodeCartEmpty,
	ErrCodeOrderNotFound,
	ErrCodePaymentFailed,
}
