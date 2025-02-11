package model

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Asset struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	Category  string    `gorm:"type:varchar(255);not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Currency  string    `gorm:"type:enum('USD', 'EUR', 'CNY');not null"`
	Amount    float64   `gorm:"type:decimal(15,2);not null"`
	CreatorID uuid.UUID
	Creator   User
	CreatedAt time.Time
	UpdatedAt time.Time
}

func InitDB() *sql.DB {
	db, err := gorm.Open("mysql", "asset:asset1234@tcp(assets-db:3306)/assets?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("failed to connect database: %v\n", err)
		panic("failed to connect database")
	}
	db.AutoMigrate(&Asset{}, &User{})
	return db.DB()
}
