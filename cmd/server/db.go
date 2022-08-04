package main

import (
	"fmt"
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
func creatOrder(db *gorm.DB, user User, product Product) error {

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
		Errno = 1
		return fmt.Errorf("create order failed: %v", Errno)
	}
	// reduce quantity
	product.Quantity--
	// update product
	err := tx.Save(&product).Error
	if err != nil {
		Errno = 2
		return fmt.Errorf("create order failed: %v", Errno)
	}

	// create order
	order := Order{UserID: user.ID, ProductId: product.ID, Status: "pending"}
	err = tx.Create(&order).Error
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
