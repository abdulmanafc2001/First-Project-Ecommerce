package models

type Wishlist struct {
	Wishtlist_ID uint `json:"wistlistid" gorm:"primaryKey;unique"`
	Product_ID  uint `json:"productid" gorm:"not null"`
	User_ID     uint `json:"userid" gorm:"not null"`
}
