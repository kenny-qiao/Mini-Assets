package repository

import (
	"database/sql"
	"fmt"

	"mini-assets/internal/model"
	"mini-assets/internal/util"

	"github.com/google/uuid"
)

// UserRepository 定义了用户相关的数据库操作接口
type UserRepository interface {
	Create(user *model.User) error
	GetByUsername(username string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uuid.UUID) error
}

// userRepository 是UserRepository接口的实现
type userRepository struct {
	db *sql.DB
}

// NewUserRepository 创建并返回一个UserRepository实例
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 插入一个新的用户记录到数据库
func (r *userRepository) Create(user *model.User) error {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `INSERT INTO users (id, username, password, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = r.db.Exec(query, uuid.New().String(), user.Username, hashedPassword, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByUsername 根据用户名从数据库中获取用户信息
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	query := `SELECT id, username, password, email, created_at, updated_at FROM users WHERE username = ?`
	row := r.db.QueryRow(query, username)

	var user model.User
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 未找到记录
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

// Update 更新数据库中的用户信息
func (r *userRepository) Update(user *model.User) error {
	// 假设不更新密码，如果需要更新密码，则需要重新加密
	query := `UPDATE users SET username = ?, email = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.Username, user.Email, user.UpdatedAt, user.ID.String())
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete 根据用户ID删除数据库中的用户记录
func (r *userRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
