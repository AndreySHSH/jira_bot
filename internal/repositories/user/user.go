package user

import "gorm.io/gorm"

type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Repository struct {
	db *gorm.DB
}

func (r Repository) CreateUser(id int64, login string, password string) {
	var u User

	u.ID = id
	u.Login = login
	u.Password = password

	r.db.Table("users").Create(&u)
}

func (r Repository) UpdateUser(id int64, login string, password string) {
	var u User

	u.ID = id
	u.Login = login
	u.Password = password

	r.db.Table("users").Updates(u)
}

func (r Repository) CheckUser(login string) bool {
	var u User
	r.db.Table("users").Where("login", login).Find(&u)

	if u.ID != 0 {
		return true
	}
	return false
}

func (r Repository) GetUserByID(id int64) *User {
	var u User
	r.db.Table("users").Where("id", id).Find(&u)

	return &u
}

func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
