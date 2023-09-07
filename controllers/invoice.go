package controllers

import (
	"fmt"
	"strconv"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

// GenetatePdf generates a PDF invoice for a specific order and user.
// @Summary Generate a PDF invoice
// @Description Generates a PDF invoice for a specific order and user, and saves it as "public/invoice.pdf".
// @Tags orders
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param order_id query string true "Order ID for which the invoice should be generated"
// @Security ApiKeyAuth
// @Success 200 {files} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/createinvoice [get]
func GenetatePdf(c *gin.Context) {
	order := c.Query("order_id")
	order_id, _ := strconv.Atoi(order)
	use, _ := c.Get("user")
	user_id := use.(models.User).User_ID

	var user models.User
	database.DB.First(&user, user_id)

	var orderdetails models.OrderItem
	err := database.DB.Where("order_id=? AND user_id=?", order_id, user_id).First(&orderdetails).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	var address models.Address
	err = database.DB.First(&address, orderdetails.Address_ID).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}

	// addresses := fmt.Sprintf("%v, %v, %v , %v", address.Building_Name, address.City, address.State, address.Zip_code)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(10, 10, 10)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 20, "E-Commerce")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(20, 10, "Invoice")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 12)
	pdf.Ln(-1)

	pdf.Cell(10, 10, "Order_id : "+order)
	pdf.Ln(-1)

	x, y := 10.0, 50.0
	lineHeight := 10.0
	colWidth := 27.0

	// Set header font style
	pdf.SetFont("Arial", "B", 12)

	// Create table headers
	pdf.Rect(x, y, colWidth, lineHeight, "D")
	pdf.CellFormat(colWidth, lineHeight, "No", "1", 0, "C", false, 0, "")
	x += colWidth
	pdf.Rect(x, y, colWidth*2, lineHeight, "D")
	pdf.CellFormat(colWidth*2, lineHeight, "Product", "1", 0, "C", false, 0, "")
	x += colWidth * 2
	pdf.Rect(x, y, colWidth, lineHeight, "D")
	pdf.CellFormat(colWidth, lineHeight, "Price", "1", 0, "C", false, 0, "")
	x += colWidth
	pdf.Rect(x, y, colWidth, lineHeight, "D")
	pdf.CellFormat(colWidth, lineHeight, "Discount", "1", 0, "C", false, 0, "")
	x += colWidth
	pdf.Rect(x, y, colWidth, lineHeight, "D")
	pdf.CellFormat(colWidth, lineHeight, "Quantity", "1", 0, "C", false, 0, "")
	x += colWidth
	pdf.Rect(x, y, colWidth, lineHeight, "D")
	pdf.CellFormat(colWidth, lineHeight, "Total Price", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	type orderdata struct {
		Product_Name string
		Price        uint
		Discount     uint
		Quantity     uint
		Total_Price  uint
	}
	var orderdatas []orderdata
	err = database.DB.Table("order_items").Select("products.product_name,order_items.price,order_items.discount,order_items.quantity,order_items.total_price").
		Joins("INNER JOIN products ON products.id=order_items.product_id").
		Where("order_items.order_id=?", order_id).Scan(&orderdatas).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	pdf.SetFont("Arial", "", 12)

	total := 0

	y += lineHeight
	x = 10.0
	for i, val := range orderdatas {
		total += int(val.Total_Price)
		index := strconv.Itoa(i + 1)
		price := strconv.Itoa(int(val.Price))
		discount := strconv.Itoa(int(val.Discount))
		quantity := strconv.Itoa(int(val.Quantity))
		totalprice := strconv.Itoa(int(val.Total_Price))

		pdf.Rect(x, y, colWidth, lineHeight, "D")
		pdf.CellFormat(colWidth, lineHeight, index, "1", 0, "C", false, 0, "")
		x += colWidth
		pdf.Rect(x, y, colWidth*2, lineHeight, "D")
		pdf.CellFormat(colWidth*2, lineHeight, val.Product_Name, "1", 0, "C", false, 0, "")
		x += colWidth * 2
		pdf.Rect(x, y, colWidth, lineHeight, "D")
		pdf.CellFormat(colWidth, lineHeight, price, "1", 0, "C", false, 0, "")
		x += colWidth
		pdf.Rect(x, y, colWidth, lineHeight, "D")
		pdf.CellFormat(colWidth, lineHeight, discount, "1", 0, "C", false, 0, "")
		x += colWidth
		pdf.Rect(x, y, colWidth, lineHeight, "D")
		pdf.CellFormat(colWidth, lineHeight, quantity, "1", 0, "C", false, 0, "")
		x += colWidth
		pdf.Rect(x, y, colWidth, lineHeight, "D")
		pdf.CellFormat(colWidth, lineHeight, totalprice, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)

	}
	pdf.Ln(-1)
	totals := strconv.Itoa(total)

	pdf.Cell(10, 10, "Total = "+totals)
	pdf.Ln(-1)
	addresses := fmt.Sprintf("%v, %v, %v , %v", address.Building_Name, address.City, address.State, address.Zip_code)
	pdf.Cell(10, 10, "Name: "+user.First_Name+" "+user.Last_Name)
	pdf.Ln(-1)

	pdf.Cell(10, 10, "Address: "+addresses)
	pdf.Ln(50)

	pdf.Cell(10, 10, "Thank you for purchasing from our website")

	pdf.OutputFileAndClose("public/invoice.pdf")
	c.HTML(200, "invoice.html", gin.H{})
}

// DownloadInvoice allows users to download the previously generated invoice in PDF format.
// @Summary Download the PDF invoice
// @Description Allows users to download the previously generated PDF invoice.
// @Tags invoice
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Security ApiKeyAuth
// @Success 200 {file} pdf "PDF invoice file for download"
// @Failure 404 "Invoice not found"
// @Router /user/downloadinvoice [get]
func DownloadInvoice(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	c.Header("Content-Type", "application/pdf")
	c.File("./public/invoice.pdf")
}
