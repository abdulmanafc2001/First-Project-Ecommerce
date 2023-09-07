package models

import "time"

type Order struct {
	Order_ID    uint   `json:"orderid" gorm:"primaryKey;unique"`
	User_ID     uint   `json:"userid" gorm:"not null"`
	Address_ID  uint   `json:"addressid" gorm:"not null"`
	Total_Price uint   `json:"totalprice" gorm:"not null"`
	Payment_ID  uint   `json:"paymentid" gorm:"not null"`
	Status      string `json:"status" gorm:"not null"`
}

type OrderItem struct {
	Order_ItemID uint      `json:"orderitemid" gorm:"primaryKey;unique"`
	User_ID      uint      `json:"userid" gorm:"not null"`
	Order_ID     uint      `json:"orderid" gorm:"not null"`
	Product_ID   uint      `json:"productid" gorm:"not null"`
	Address_ID   uint      `json:"addressid" gorm:"not null"`
	Catagory     string    `json:"catagory" gorm:"not null"`
	Brand        string    `json:"brand" gorm:"not null"`
	Quantity     uint      `json:"quantity" gorm:"not null"`
	Price        uint      `json:"price" gorm:"not null"`
	Total_Price  uint      `json:"totalprice" gorm:"not null"`
	Discount     uint      `json:"discount"`
	Cart_ID      uint      `json:"cartid" gorm:"not null"`
	Status       string    `json:"status" gorm:"not null"`
	Created_at   time.Time `json:"createdat"`
}

type Payment struct {
	Payment_ID     uint      `json:"paymentid" gorm:"primaryKey;unique"`
	Payment_Type   string    `json:"paymenttype" gorm:"not null"`
	Total_Amount   uint      `json:"totalamount" gorm:"not null"`
	Payment_Status string    `json:"paymentstatus" gorm:"not null"`
	User_ID        uint      `json:"userid" gorm:"not null"`
	Date           time.Time `json:"date"`
}

type RazorPay struct {
	User_id          uint   `json:"userid"`
	RazorPayment_id  string `json:"razorpaymentid" gorm:"primaryKey"`
	RazorPayOrder_id string `json:"razorpayorderid"`
	Signature        string `json:"signature"`
	AmountPaid       string `json:"amountpaid"`
}
