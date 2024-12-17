package domain

// User 领域对象User，是 DDD 中的聚合根
type User struct {
	Id       int64
	Email    string
	Phone    string
	Password string
	CreateAt int64
}
