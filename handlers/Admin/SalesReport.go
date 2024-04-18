package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/tealeg/xlsx"
)

func SalesReport() {
	var order []models.Order

	var totalAmount float64
	if err := initializers.DB.Find(&order); err.Error != nil {
		for _, val := range order {
			totalAmount += float64(val.OrderPrice)
		}
	}
}

// @Summary Generate Excel file  with sales report data
// @Description  This can only be used by admin to  generate excel file of the sales report
// @Tags Admin-SalesReport
// @Accept json
// @Produce json
// @Success 200 {json} json	"Excel file generated and sent successfully"
// @Failure 401 {json} json "Failed to Generate  excel file"
// @Router /admin/SalesReport/downloadexcel [get]
func GenerateSalesReport(c *gin.Context) {

	var OrderData []models.OrderItem
	if err := initializers.DB.Preload("Product").Preload("Order.User").Find(&OrderData).Error; err != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to fetch sales data",
			"status": 401,
		})
		return
	}
	var cancelCount int
	for _, view := range OrderData {
		if view.Orderstatus == "Order cancelled" {
			cancelCount++
		}
	}
	//========================= Create excel file =======================
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sales Report")
	if err != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to create Excel sheet",
			"status": 401,
		})
		return
	}
	//===================== Add headers to the excel sheet ===========================
	headers := []string{"Order ID", "Customer Name", "Product Name", "Order Date", "Total Amount", "Order Status"}
	row := sheet.AddRow()
	for _, header := range headers {
		cell := row.AddCell()
		cell.Value = header
	}
	//===================== Add sales data =======================
	for _, sale := range OrderData {
		row := sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(int(sale.OrderID))
		row.AddCell().Value = sale.Order.User.Name
		row.AddCell().Value = sale.Product.Product_Name
		row.AddCell().Value = sale.Order.OrderDate.Format("2016-02-01")
		row.AddCell().Value = fmt.Sprintf("%d", sale.Order.OrderPrice)
		row.AddCell().Value = sale.Orderstatus
	}
	// row.AddCell().Value = strconv.Itoa(cancelCount)
	//========================= Save Excel File ============================
	path := "/home/home/Brototype/Brototype/brocamp/week-10/sales_report.xlsx"
	if err := file.Save(path); err != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to fetch sales data",
			"status": 401,
		})
		return
	}
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path))
	c.Writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(path)

	c.JSON(200, gin.H{
		"message": "Excel file generated and downloaded successfully",
		"status":  200,
	})
	fmt.Println("Excel file generated and sent successfully")
}

// @Summary Generate PDF file  with sales report data
// @Description  This can only be used by admin to  generate pdf file of the sales report
// @Tags Admin-SalesReport
// @Accept json
// @Produce json
// @Success 200 {json} json	"PDF file generated and sent successfully"
// @Failure 401 {json} json "Failed to Generate  pdf file"
// @Router /admin/SalesReport/downloadpdf [get]
func SalesReportPDF(c *gin.Context) {
	var OrderData []models.OrderItem
	if err := initializers.DB.Preload("Product").Preload("Order.User").Find(&OrderData).Error; err != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to fetch sales data",
			"status": 401,
		})
		return
	}

	//================== Creating  new PDF document ========================
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 12)

	//================== Add headers to the PDF ============================
	headers := []string{"Order ID", "Customer", "Product", "Order Date", "Total Amount", "Order Status"}
	for _, header := range headers {
		pdf.Cell(40, 10, header)
	}
	pdf.Ln(-1)

	//===================== Add sales data to the PDF =======================

	for _, sale := range OrderData {
		pdf.Cell(40, 10, strconv.Itoa(int(sale.OrderID)))
		pdf.Cell(40, 10, sale.Order.User.Name)
		pdf.Cell(40, 10, sale.Product.Product_Name)
		pdf.Cell(40, 10, sale.Order.OrderDate.Format("2016-02-01"))
		pdf.Cell(40, 10, fmt.Sprintf("%.2f", sale.Subtotal))
		pdf.Cell(40, 10, sale.Orderstatus)
		pdf.Ln(-1)
	}

	//===================== Save PDF file ===================================
	pdfPath := "/home/home/Brototype/Brototype/brocamp/week-10/sales_report.pdf"
	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to generate PDF file",
			"status": 401,
			"result": err,
		})
		return
	}
	//====================== PDF file download ==============================
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", pdfPath))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(pdfPath)

	c.JSON(200, gin.H{
		"message": "PDF file generated and downloaded successfully",
		"status":  200,
	})
	fmt.Println("PDF file generated and sent successfully")
}
