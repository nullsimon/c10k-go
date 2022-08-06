package main

import (
	"context"
	"fmt"
	redis2 "github.com/go-redis/redis/v9"
	"github.com/nullsimon/c10k-go/cmd/server/redis"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code     string
	Price    uint
	Quantity uint
}

type User struct {
	gorm.Model
	Name string
}

type Order struct {
	gorm.Model
	UserID    uint   `gorm:"index"`
	ProductId uint   `gorm:"index"`
	Status    string `gorm:"not null;default:null"`
}

// success
func creatOrder(db *gorm.DB, redisClient *redis2.Client, user User, product Product) error {

	// begin transaction
	tx := db.Begin()
	Errno := 0

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("create order failed: %v\n", err)
			tx.Rollback()
		} else {
			if Errno != 0 {
				fmt.Printf("create order failed: %v\n", Errno)
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}()

	// redis descrease quantity
	ok := redis.DecreaseQuantity(context.Background(), redisClient, product.Code, 1)
	if !ok {
		Errno = 2
		return fmt.Errorf("redis decrease quantity failed: %v", Errno)
	}

	// create order
	order := Order{UserID: user.ID, ProductId: product.ID, Status: "pending"}
	err := tx.Create(&order).Error
	if err != nil {
		Errno = 3
		return fmt.Errorf("create order failed: %v", Errno)
	}
	return nil
}

// fail
func creatOrderFail(db *gorm.DB, user User, product Product) {

	// begin transaction
	tx := db.Begin()
	Errno := 0

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("create order failed: %v\n", err)
			tx.Rollback()
		} else {
			if Errno != 0 {
				fmt.Printf("create order failed: %v\n", Errno)
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}()

	// lock read product
	tx.Set("gorm:query_option", "FOR UPDATE").First(&product)

	// check quantity
	if product.Quantity < 1 {
		return
	}
	// reduce quantity
	product.Quantity--
	// update product
	err := tx.Save(&product).Error
	if err != nil {
		return
	}
	order := Order{UserID: user.ID, ProductId: product.ID}
	err = tx.Create(&order).Error
	if err != nil {
		fmt.Printf("create order failed: %v\n", err)
		Errno = 1
		return
	}
}
