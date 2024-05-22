package models

import "time"

type OTPS struct {
	ID        int64     `json:"Id,omitempty" bson:"Id,omitempty"`
	Otp       string    `json:"otp,omitempty" bson:"otp,omitempty" validate:"required"`
	UserID    int64     `json:"userId,omitempty" bson:"userId,omitempty" validate:"required"`
	WalletID  int64     `json:"walletId,omitempty" bson:"walletId,omitempty" validate:"required"`
	Status    string    `json:"status,omitempty" bson:"status,omitempty" validate:"required"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty" validate:"required"`
}
