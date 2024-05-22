package models

type Wallet struct {
	ID      int64  `json:"Id,omitempty" bson:"Id,omitempty"`
	Balance float64 `json:"balance,omitempty" bson:"balance,omitempty" validate:"required"`
	UserID  int64  `json:"userId,omitempty" bson:"userId,omitempty" validate:"required"`
	Status  string `json:"status,omitempty" bson:"status,omitempty" validate:"required"`
}
