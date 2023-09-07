package controllers

import (
	"fmt"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/tealeg/xlsx"
)

// SalesReport generates a sales report within a specified date range and exports it to an Excel file.
// @Summary Generate Sales Report
// @Description Generates a sales report within a specified date range and exports it to an Excel file.
// @Tags reports
// @Accept json
// @Produce json
// @Param startingdate query string true "Starting date in the format YYYY-MM-DD"
// @Param endingdate query string true "Ending date in the format YYYY-MM-DD"
// @Success 200 {html} HTML "Sales report generated successfully"
// @Failure 400 {string} ErrorResponse
// @Router /admin/salesreport [get]
func SalesReport(c *gin.Context) {
	startingDate := c.Query("startingdate")
	endingDate := c.Query("endingdate")

	start, err := time.Parse("2006-01-02", startingDate)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Starting date conversion error",
		})
		return
	}
	end, err := time.Parse("2006-01-02", endingDate)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Ending date conversion error",
		})
		return

	}
	type orderDetails struct {
		Product_Name string
		Quantity     string
		Price        uint
		Total_Price  uint
		Status       string
		Order_ID     uint
	}
	//Fetching data from database and inner joins product table knowing product details
	var orders []orderDetails
	err = database.DB.Table("order_items").
		Select("products.product_name,order_items.quantity,order_items.price,order_items.total_price,order_items.status,order_items.order_id").
		Joins("INNER JOIN products ON products.id=order_items.product_id").Where("order_items.created_at BETWEEN ? AND ?", start, end).
		Scan(&orders).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Scanning error",
		})
		return
	}
	fmt.Println(orders)
	f := excelize.NewFile()
	sheet := "sheet1"
	intex := f.NewSheet(sheet)

	f.SetCellValue(sheet, "A1", "Product name")
	f.SetCellValue(sheet, "B1", "Quantity")
	f.SetCellValue(sheet, "C1", "Price")
	f.SetCellValue(sheet, "D1", "Total price")
	f.SetCellValue(sheet, "E1", "Status")
	f.SetCellValue(sheet, "F1", "Order_id")

	for i, v := range orders {
		k := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", k), v.Product_Name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", k), v.Quantity)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", k), v.Price)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", k), v.Total_Price)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", k), v.Status)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", k), v.Order_ID)
	}
	f.SetActiveSheet(intex)

	if err := f.SaveAs("./public/salesreport.xlsx"); err != nil {
		fmt.Println(err)
		return
	}
	convertintoPdf(c)
	c.HTML(200, "salesreport.html", gin.H{})
}

func convertintoPdf(c *gin.Context) {
	file, err := xlsx.OpenFile("./public/salesreport.xlsx")

	if err != nil {
		fmt.Println("File open error")
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(20, 10, 20)
	pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(10, 10, "Sales report")
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 8)
	// Convertig each cell in the Excel file to a PDF cell
	for _, sheet := range file.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				//if there is any empty cell values skiping that
				if cell.Value == "" {
					continue
				}

				pdf.CellFormat(25, 10, cell.Value, "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)
		}
	}
	err = pdf.OutputFileAndClose("./public/salesreport.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("PDF saved successfully.")
}

// DownloadExel allows users to download the sales report Excel file.
// @Summary Download Sales Report Excel File
// @Description Allows users to download the sales report Excel file as an attachment.
// @Tags reports
// @Produce json
// @Success 200 {file} file "salesreport.xlsx"
// @Failure 404 {json} JSON "File not found" (when the Excel file is not found)
// @Failure 500 {json} JSON "Internal server error" (for other errors)
// @Router /admin/salesreport/xlsx [get]
func DownloadExel(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=salesreport.xlsx")
	c.Header("Content-Type", "application/xlsx")
	c.File("./public/salesreport.xlsx")
}

// Downloadpdf allows users to download the sales report PDF file.
// @Summary Download Sales Report PDF File
// @Description Allows users to download the sales report PDF file as an attachment.
// @Tags reports
// @Produce json
// @Success 200 {file} file "salesreport.pdf"
// @Failure 404 {json} JSON "File not found" (when the PDF file is not found)
// @Failure 500 {json} JSON "Internal server error" (for other errors)
// @Router /admin/salesreport/pdf [get]
func Downloadpdf(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=salesreport.pdf")
	c.Header("Content-Type", "application/pdf")
	c.File("./public/salesreport.pdf")
}
