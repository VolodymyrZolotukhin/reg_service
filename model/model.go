package model

import (
	"gorm.io/gorm"
)

type User struct {
	Id       int    `gorm:"primary key;autoIncrement" json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *User) Create(db *gorm.DB) (int, error) {
	result := db.Create(u)

	return u.Id, result.Error
}

func (u *User) GetById(db *gorm.DB, id int) error {
	result := db.First(u, "id = ?", id)

	return result.Error
}

func (u *User) GetByLogin(db *gorm.DB, login string) error {
	result := db.Find(u, "login = ?", login)

	return result.Error
}

func (u *User) Update(db *gorm.DB) error {
	err := u.GetById(db, u.Id)

	if err != nil {
		return err
	}

	result := db.Where("id = ?", u.Id).Updates(u)

	return result.Error
}
