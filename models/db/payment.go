package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	ID         int64   `json:"id" gorm:"column:id;primary_key"`
	PaymentID  string  `json:"payment_id" gorm:"column:payment_id;not null"`
	CardNumber string  `gorm:"column:card_number;not null"`
	ExpMonth   int     `gorm:"column:exp_month;not null"`
	ExpYear    int     `gorm:"column:exp_year;not null"`
	Amount     float64 `gorm:"column:amount;not null"`
	Currency   string  `gorm:"column:currency;not null"`
	CVV        int     `gorm:"column:cvv;not null"`
	Status     string  `json:"status"`
}

func (Payment) TableName() string {
	return "payments"
}

// Card represents a row in the "cards" table
type Card struct {
	gorm.Model
	CardNumber string  `gorm:"column:card_number;not null"`
	Balance    float64 `gorm:"column:balance;not null"`
}

// TableName specifies the table name for the Card struct
func (Card) TableName() string {
	return "cards"
}

// Fraud represents a row in the "frauds" table
type Fraud struct {
	gorm.Model
	CardNumber string  `gorm:"column:card_number;not null"`
	Amount     float64 `gorm:"column:amount;not null"`
	Currency   string  `gorm:"column:currency;not null"`
}

// TableName specifies the table name for the Fraud struct
func (Fraud) TableName() string {
	return "frauds"
}
