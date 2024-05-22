package models

import "time"

type Transactions struct {
	ID        int64     `json:"Id,omitempty" bson:"Id,omitempty"`
	Amount    float64   `json:"amount,omitempty" bson:"amount,omitempty" validate:"required"`
	WalletID  int64     `json:"walletId,omitempty" bson:"walletId,omitempty" validate:"required"`
	Status    string    `json:"status,omitempty" bson:"status,omitempty" validate:"required"`
	Type      string    `json:"type,omitempty" bson:"type,omitempty" validate:"required"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty" validate:"required"`
}
