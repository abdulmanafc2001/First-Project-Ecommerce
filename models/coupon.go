package models

import "time"

type Coupon struct {
	CouponId      uint      `json:"couponid" gorm:"primaryKey;unique"`
	Coupon_Code   string    `json:"couponcode" gorm:"not null"`
	Starting_Time time.Time `json:"startingtime" gorm:"not null"`
	Ending_Time   time.Time `json:"endingtime" gorm:"not null"`
	Value         uint      `json:"value" gorm:"not null"`
	Type          string    `json:"type" gorm:"not null"`
	Max_Discount  uint      `json:"maxdiscount" gorm:"not null"`
	Min_Discount  uint      `json:"mindiscount" gorm:"not null"`
	Cancel        bool      `json:"cancel" gorm:"default:false"`
}
