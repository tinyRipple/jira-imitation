package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var ErrorUserEmailDuplicated = errors.New("ErrorUserEmailDuplicated")

type UserDao struct {
	db *gorm.DB
}

func InitUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

type User struct {
	Id          int64  `gorm:"primaryKey,autoIncrement"`
	Email       string `gorm:"unique"`
	Password    string
	CreatedTime int64
	UpdatedTime int64
	DeletedTime int64
}

func (dao *UserDao) Insert(ctx context.Context, user User) error {
	err := dao.db.WithContext(ctx).Create(&user).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == 1062 {
			return ErrorUserEmailDuplicated
		}
	}
	return nil
}