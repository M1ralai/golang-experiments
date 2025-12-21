package domain

type UserRepository interface {
	GetAll() ([]User, error)

	GetByUsername(username string) (*User, error)

	Create(user *User) error

	Delete(id int) error
}
