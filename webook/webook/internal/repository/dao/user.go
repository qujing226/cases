package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("duplicate email")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

// UserDAOer 定义了用户数据访问对象的接口
type UserDAOer interface {
	Insert(ctx context.Context, u User) error
	FindById(ctx context.Context, id int64) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
}

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAOer {
	return &UserDAO{
		db: db,
	}
}

type User struct {
	Id        int64          `gorm:"primaryKey,autoIncrement"`
	Email     sql.NullString `gorm:"unique"`
	Password  string
	Phone     sql.NullString `gorm:"unique"`
	CreatedAt int64
	UpdatedAt int64
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreatedAt = now
	u.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const uniqueViolation = 1062
		if me.Number == uniqueViolation {
			return ErrUserDuplicate
		}
	}
	return err
}
func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}
func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, ErrUserNotFound
	}
	return u, err
}
