package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository"
)

var (
	ErrUserNotFound          = repository.ErrUserNotFound
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("nickname or password wrong")
)

type UserServicer interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Profile(ctx context.Context, u domain.User) (domain.User, error)
}

type UserService struct {
	repo repository.UserRepositoryer
}

func NewUserService(repo repository.UserRepositoryer) UserServicer {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 考虑加密放在哪里 + 存放
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	err = svc.repo.Create(ctx, u)
	if err == ErrUserDuplicate {
		return ErrUserDuplicate
	}

	return err
}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return u, ErrInvalidUserOrPassword
		}
		return u, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// 打日志
		fmt.Println(err)
		return u, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) Profile(ctx context.Context, u domain.User) (domain.User, error) {
	user, err := svc.repo.FindById(ctx, u.Id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err == nil {
		return u, err
	}
	if err != ErrUserNotFound {
		return u, err
	}
	// ErrUserNotFound 用户不存在
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil {
		if err != repository.ErrUserDuplicate {
			return u, err
		}
		return u, ErrUserDuplicate
	}
	return svc.repo.FindByPhone(ctx, phone)
}
