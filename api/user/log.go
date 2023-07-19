package user

const (
	OperationRegister = "register"
	OperationLogin    = "login"
	OperationPurchase = "purchase"
	// 添加其他操作类型...
)

// 日志结构
type UserLog struct {
	UserID    uint32 `db:"user_id"`
	IP        string `db:"ip"`
	Operation string `db:"operation"`
	InitTime  int64  `db:"init_time"`
	Details   string `db:"details"`
}
