package pdf

import (
	"LinksChecker/internal/models"
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// GenerateReport creates a PDF report from tasks
func GenerateReport(tasks []*models.Task) []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()
	addHeader(pdf, len(tasks))
	addContent(pdf, tasks)

	buffer := new(bytes.Buffer)
	pdf.Output(buffer)

	return buffer.Bytes()
}

func addHeader(pdf *gofpdf.Fpdf, taskCount int) {
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Link Status Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(8)
	pdf.Cell(0, 6, fmt.Sprintf("Total tasks: %d", taskCount))
	pdf.Ln(12)
}

func addContent(pdf *gofpdf.Fpdf, tasks []*models.Task) {
	pdf.SetFont("Arial", "", 12)

	for _, task := range tasks {
		addTaskSection(pdf, task)
	}
}

func addTaskSection(pdf *gofpdf.Fpdf, task *models.Task) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Task #%d", task.ID))
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	for link, status := range task.Links {
		statusSymbol := getStatusSymbol(status)
		pdf.Cell(0, 6, fmt.Sprintf("%s %s - %s", statusSymbol, link, status))
		pdf.Ln(6)
	}
	pdf.Ln(5)
}

func getStatusSymbol(status string) string {
	if status == "available" {
		return "✓"
	}
	return "✗"
}
