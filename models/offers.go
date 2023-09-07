package models

type Catagory_Offer struct {
	ID          uint `json:"id" gorm:"primaryKey;unique"`
	Catagory_Id uint `json:"catagoryid" gorm:"not null"`
	Offer       bool `json:"offer" gorm:"not null"`
	Percentage  uint `json:"percentage"`
}
