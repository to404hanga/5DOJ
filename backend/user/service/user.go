package service

import (
	"5DOJ/user/domain"
	"5DOJ/user/global"
	"5DOJ/user/model"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	defaultPassword []byte
}

var _ IUserService = (*UserService)(nil)

func NewUserService(defaultPassword string) *UserService {
	return &UserService{
		defaultPassword: []byte(defaultPassword),
	}
}

func (u *UserService) SignUp(ctx context.Context, userView domain.UserView) (err error) {
	bytes, _ := bcrypt.GenerateFromPassword(u.defaultPassword, bcrypt.DefaultCost)
	user := model.User{
		Uid:             userView.Uid,
		Name:            userView.Name,
		Password:        string(bytes),
		TelephoneNumber: userView.TelephoneNumber,
		Gender:          userView.Gender,
	}

	return global.MySQL.WithContext(ctx).Create(&user).Error
}

func (u *UserService) Login(ctx context.Context, uid, password string) (userView domain.UserView, err error) {
	var user model.User
	err = global.MySQL.WithContext(ctx).Where("uid = ?", uid).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.UserView{}, errors.New("账号或密码错误")
		}
		return domain.UserView{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.UserView{}, errors.New("账号或密码错误")
	}

	userView = domain.UserView{
		Uid:             user.Uid,
		Name:            user.Name,
		TelephoneNumber: user.TelephoneNumber,
		Gender:          user.Gender,
	}

	return userView, nil
}

func (u *UserService) UpdatePassword(ctx context.Context, uid, password, confirmPassword string) (err error) {
	if password != confirmPassword {
		return errors.New("两次密码不一致")
	}

	tx := global.MySQL.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("%v", r)
		}
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var hashedPassword string
	err = tx.WithContext(ctx).Model(&model.User{}).Where("uid = ?", uid).Select("password").Scan(&hashedPassword).Error
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("原密码错误")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = tx.WithContext(ctx).Model(&model.User{}).Where("uid = ?", uid).Update("password", string(bytes)).Error
	if err != nil {
		return err
	}

	return
}

func (u *UserService) GetInfoByUid(ctx context.Context, uid string) (userView domain.UserView, err error) {
	var user model.User
	err = global.MySQL.WithContext(ctx).Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return domain.UserView{}, err
	}

	userView = domain.UserView{
		Uid:             user.Uid,
		Name:            user.Name,
		TelephoneNumber: user.TelephoneNumber,
		Gender:          user.Gender,
	}

	return
}
