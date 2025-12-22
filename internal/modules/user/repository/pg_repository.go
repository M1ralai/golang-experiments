package repository

import (
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/user/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) domain.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetAll() ([]domain.User, error) {
	users := []domain.User{}
	query := `SELECT id, username, '' as password, role, COALESCE(ad, '') as ad, COALESCE(soyad, '') as soyad, COALESCE(telefon, '') as telefon, COALESCE(email, '') as email FROM users ORDER BY id ASC`
	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) GetByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, username, password, role, COALESCE(ad, '') as ad, COALESCE(soyad, '') as soyad, COALESCE(telefon, '') as telefon, COALESCE(email, '') as email FROM users WHERE username = $1`
	err := r.db.Get(user, query, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (username, password, role, ad, soyad, telefon, email) VALUES (:username, :password, :role, :ad, :soyad, :telefon, :email)`
	_, err := r.db.NamedExec(query, user)
	return err
}

func (r *PostgresUserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
