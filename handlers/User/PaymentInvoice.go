package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

// @Summary Generate payment invoice for a specific order
// @Description Generate a payment invoice for a specific order and download it as a PDF
// @Tags User-Invoice
// @Param ID path string true "Order ID"
// @Accept json
// @Produce application/pdf
// @Success 200 {json} json "PDF file generated and downloaded successfully"
// @Failure 500 {json} json "Failed to fetch invoice data or generate PDF file"
// @Router /invoice/{ID} [get]
func PymentInvoice(c *gin.Context) {
	var orderItems []models.OrderItem
	orderID := c.Param("ID")
	if err := initializers.DB.Where("order_id = ? AND orderstatus = ?", orderID, "delivered").Preload("Product").Preload("Order.Address").Find(&orderItems); err.Error != nil {
		c.JSON(401, gin.H{
			"Error":  "Failed to fetch invoice data",
			"status": 401,
		})
		return
	}
	//=========== Generate Invoide ID ================
	const codeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, 4)
	for i := range code {
		code[i] = codeCharset[rand.Intn(len(codeCharset))]
	}
	InvoiceID := "INV-" + string(code) + orderID
	fmt.Println("code------->", InvoiceID)
	//========== Adding A4 size pge ===================
	marginX := 10.0
	marginY := 20.0
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginX, marginY, marginX)
	pdf.AddPage()

	// ============ Logo of the company =================
	pdf.ImageOptions("images/watch1 img2.jpg", 11, 6, 25, 25, false, gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")
	pdf.Ln(5)

	// ============= Company Name and Billing Address =======================
	Date := ""
	for _, value := range orderItems {
		Date = value.Order.OrderDate.Format("2006-01-02")
	}
	pdf.SetFont("Arial", "B", 32)
	_, lineHeight := pdf.GetFontSize()
	currentY := pdf.GetY() + lineHeight
	// lineBreak := lineHeight + float64(1)
	pdf.SetY(currentY)
	pdf.CellFormat(40, 10, "Ecom", "0", 0, "L", false, 0, "")
	pdf.Ln(0)

	pdf.SetFont("Arial", "B", 32)
	_, lineHeight = pdf.GetFontSize()
	pdf.SetXY(100, currentY-lineHeight)
	pdf.CellFormat(40, -20, "INVOICE", "0", 0, "R", false, 0, "")
	pdf.Ln(30)

	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(109, -40, "Invoice No:", "0", 0, "R", false, 0, "")
	pdf.Ln(0)
	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(140, -40, InvoiceID, "0", 0, "R", false, 0, "")
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(109, -40, "Order Date:", "0", 0, "R", false, 0, "")
	pdf.Ln(0)
	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(133, -40, Date, "0", 0, "R", false, 0, "")
	pdf.Ln(0)

	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(40, 8, "Bill To:")
	pdf.Ln(5)
	//================ Billing Address ======================
	pdf.SetFont("Arial", "", 12)
	for _, address := range orderItems {
		pdf.CellFormat(40, 10, address.Order.Address.Address, "0", 0, "L", false, 0, "")
		pdf.Ln(5)
		pdf.CellFormat(40, 10, address.Order.Address.City, "0", 0, "L", false, 0, "")
		pdf.Ln(5)
		pdf.CellFormat(40, 10, address.Order.Address.Country, "0", 0, "L", false, 0, "")
		pdf.Ln(5)
		pdf.CellFormat(40, 10, address.Order.Address.State, "0", 0, "L", false, 0, "")
		pdf.Ln(5)
		pdf.CellFormat(40, 10, strconv.Itoa(address.Order.Address.Pincode), "0", 0, "L", false, 0, "")
		pdf.Ln(5)
		pdf.CellFormat(10, 10, "Tel:", "0", 0, "L", false, 0, "")
		pdf.CellFormat(40, 10, strconv.Itoa(address.Order.Address.Phone), "0", 0, "L", false, 0, "")
		break
	}
	pdf.Ln(-1)

	//==================== Add headers to the PDF =============
	pdf.SetFont("Arial", "B", 16)
	lineHt := 10.0
	const colNumber = 6
	headers := [colNumber]string{"No", "Name", "Discritption", "Qty", "Unit Price", "Price"}
	colWidth := [colNumber]float64{13.0, 30.0, 70.0, 15.0, 30.0, 25.0}
	pdf.SetFillColor(200, 200, 200)
	for colJ := 0; colJ < colNumber; colJ++ {
		pdf.CellFormat(colWidth[colJ], lineHt, headers[colJ], "1", 0, "CM", true, 0, "")
	}
	pdf.Ln(-1)

	//==================== Add datas to the PDF ===================
	pdf.SetLeftMargin(10)
	pdf.SetFont("Arial", "", 14)
	count := 0
	var GrandTotal = 0
	var Discount = 0
	DeliveryCharge := 0
	for _, invoice := range orderItems {
		count++
		DeliveryCharge = invoice.Order.DeliveryCharge
		GrandTotal += int(invoice.Subtotal)
		Discount = GrandTotal - (invoice.Order.OrderPrice - DeliveryCharge)
		pdf.CellFormat(colWidth[0], lineHt, strconv.Itoa(count), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth[1], lineHt, invoice.Product.Product_Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(colWidth[2], lineHt, invoice.Product.Description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(colWidth[3], lineHt, strconv.Itoa(invoice.OrderQuantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%.2d", invoice.Product.Price), "1", 0, "L", false, 0, "")
		pdf.CellFormat(colWidth[5], lineHt, fmt.Sprintf("%.2f", invoice.Subtotal), "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.SetFontStyle("B")
	leftIndent := 15.0
	for i := 0; i < 3; i++ {
		leftIndent += colWidth[i]
	}
	Total := (GrandTotal + DeliveryCharge) - Discount
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[4], lineHt, "Subtotal", "1", 0, "L", false, 0, "")
	pdf.CellFormat(colWidth[5], lineHt, fmt.Sprintf("%.2d", GrandTotal), "1", 0, "L", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[4], lineHt, "Discount", "1", 0, "L", false, 0, "")
	pdf.CellFormat(colWidth[5], lineHt, fmt.Sprintf("%.2d", Discount), "1", 0, "L", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[4], lineHt, "Deli.Charge", "1", 0, "L", false, 0, "")
	pdf.CellFormat(colWidth[5], lineHt, fmt.Sprintf("%.2d", DeliveryCharge), "1", 0, "L", false, 0, "")
	pdf.Ln(-1)
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[4], lineHt, "Grand Total", "1", 0, "L", false, 0, "")
	pdf.CellFormat(colWidth[5], lineHt, fmt.Sprintf("%.2d", Total), "1", 0, "L", false, 0, "")

	//===================== Save PDF file ===================================
	pdfPath := "/home/home/Brototype/Brototype/brocamp/week-10/Invoice.pdf"
	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to generate PDF file",
			"status": 401,
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
