package models

type Cart struct {
	ID              uint `json:"id" gorm:"primaryKey;unique"`
	Product_ID      uint `json:"productid" gorm:"not null"`
	Quantity        uint `json:"quantity" gorm:"not null"`
	Price           uint `json:"price" gorm:"not null"`
	Total_Price     uint `json:"totalprice" gorm:"not null"`
	Catagory_ID     uint `json:"catagoryid" gorm:"not null"`
	User_ID         uint `json:"userid" gorm:"not null"`
	Coupon_Applied  bool `json:"couponapplied" gorm:"default:false"`
	Coupon_Discount uint `json:"discount"`
	Catagory_Offer  uint `json:"catagoryoffer"`
}
