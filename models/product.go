package models

import "gorm.io/gorm"

type Catagory struct {
	Catagory_ID   uint   `json:"catagoryid" gorm:"primaryKey;unique"`
	Catagory_Name string `json:"catagoryname" gorm:"not null"`
	Unlist        bool   `json:"unlist" gorm:"default:false"`
}

type Product struct {
	gorm.Model
	Product_Name string `json:"productname" gorm:"not null"`
	Descreption  string `json:"descreption" gorm:"not null"`
	Stock        uint   `json:"stock" gorm:"not null"`
	Price        uint   `json:"price" gorm:"not null"`
	Catagory_ID  uint   `json:"catagoryid" gorm:"not null"`
	Brand_ID     uint   `json:"brandid" gorm:"not null"`
}
type Brand struct {
	Brand_ID   uint   `json:"brandid" gorm:"primaryKey;unique"`
	Brand_Name string `json:"brandname" gorm:"not null"`
}
type Image struct {
	Id         uint   `json:"id" gorm:"primaryKey;unique"`
	Product_id uint   `json:"productid" gorm:"not null"`
	Image      string `json:"image" gorm:"not null"`
}
