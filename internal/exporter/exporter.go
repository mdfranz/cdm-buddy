package exporter

import (
	"encoding/csv"
	"fmt"
	"os"

	"cdmbuddy/internal/model"
	"github.com/xuri/excelize/v2"
)

func ExportToCSV(matrix model.Matrix, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Asset Class", "Instance", "Function", "Technology", "People", "Process"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, asset := range model.AssetClasses {
		instances := matrix[asset]
		for _, instance := range instances {
			for _, funcName := range model.Functions {
				cell := instance.Cells[funcName]
				// Only export cells that have at least one value
				if cell.Tech != "" || cell.People != "" || cell.Process != "" {
					row := []string{
						asset,
						instance.Name,
						funcName,
						cell.Tech,
						cell.People,
						cell.Process,
					}
					if err := writer.Write(row); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func ExportToExcel(matrix model.Matrix, path string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "Cyber Defense Matrix"
	f.SetSheetName("Sheet1", sheetName)

	// Styles
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	boomStyleBase := &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14, Color: "FFFFFF"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	}

	governBoomStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"6AA84F"}, Pattern: 1},
		Font:      boomStyleBase.Font,
		Alignment: boomStyleBase.Alignment,
		Border:    boomStyleBase.Border,
	})

	leftBoomStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4F81BD"}, Pattern: 1},
		Font:      boomStyleBase.Font,
		Alignment: boomStyleBase.Alignment,
		Border:    boomStyleBase.Border,
	})

	rightBoomStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"C0504D"}, Pattern: 1},
		Font:      boomStyleBase.Font,
		Alignment: boomStyleBase.Alignment,
		Border:    boomStyleBase.Border,
	})

	// 1. Add Boom Headers
	f.SetCellValue(sheetName, "C1", "GOVERN (Cross-Cutting)")
	f.SetCellStyle(sheetName, "C1", "C1", governBoomStyle)

	f.MergeCell(sheetName, "D1", "E1")
	f.SetCellValue(sheetName, "D1", "LEFT OF BOOM (Pre-Event)")
	f.SetCellStyle(sheetName, "D1", "E1", leftBoomStyle)

	f.MergeCell(sheetName, "F1", "H1")
	f.SetCellValue(sheetName, "F1", "RIGHT OF BOOM (Post-Event)")
	f.SetCellStyle(sheetName, "F1", "H1", rightBoomStyle)

	// 2. Add Function Headers
	for i, funcName := range model.Functions {
		col, _ := excelize.ColumnNumberToName(i + 3) // Start from C (3)
		cell := col + "2"
		f.SetCellValue(sheetName, cell, funcName)
		
		color := model.FunctionColors[funcName]
		style, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Size: 12},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
		})
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	// Add "Asset Instance" header
	f.SetCellValue(sheetName, "B2", "Asset Instance")
	f.SetCellStyle(sheetName, "B2", "B2", headerStyle)

	// 3. Add Asset Headers and Data
	currentRow := 3
	for _, asset := range model.AssetClasses {
		instances := matrix[asset]
		if len(instances) == 0 {
			continue
		}

		startRow := currentRow
		for _, instance := range instances {
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), instance.Name)
			f.SetCellStyle(sheetName, fmt.Sprintf("B%d", currentRow), fmt.Sprintf("B%d", currentRow), headerStyle)

			for c, funcName := range model.Functions {
				col, _ := excelize.ColumnNumberToName(c + 3)
				cell := fmt.Sprintf("%s%d", col, currentRow)
				
				cellData := instance.Cells[funcName]
				value := fmt.Sprintf("TECH: %s\nPEOPLE: %s\nPROCESS: %s", cellData.Tech, cellData.People, cellData.Process)
				f.SetCellValue(sheetName, cell, value)

				color := model.FunctionColors[funcName]
				
				// Check if we need thick right border for "Protect" (column E now, index 2 in Functions)
				var borders []excelize.Border
				borders = append(borders, excelize.Border{Type: "left", Color: "000000", Style: 1})
				borders = append(borders, excelize.Border{Type: "top", Color: "000000", Style: 1})
				borders = append(borders, excelize.Border{Type: "bottom", Color: "000000", Style: 1})
				
				if funcName == "Protect" {
					borders = append(borders, excelize.Border{Type: "right", Color: "000000", Style: 5}) // 5 is thick
				} else {
					borders = append(borders, excelize.Border{Type: "right", Color: "000000", Style: 1})
				}

				style, _ := f.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "top", WrapText: true},
					Fill:      excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
					Border:    borders,
				})
				f.SetCellStyle(sheetName, cell, cell, style)
			}
			currentRow++
		}
		
		endRow := currentRow - 1
		cellA := fmt.Sprintf("A%d", startRow)
		if startRow == endRow {
			f.SetCellValue(sheetName, cellA, asset)
			f.SetCellStyle(sheetName, cellA, cellA, headerStyle)
		} else {
			f.MergeCell(sheetName, cellA, fmt.Sprintf("A%d", endRow))
			f.SetCellValue(sheetName, cellA, asset)
			f.SetCellStyle(sheetName, cellA, fmt.Sprintf("A%d", endRow), headerStyle)
		}
	}

	// 4. Thick border for Boom Boundary
	// Protect is column E (5). We apply it to rows 1 and 2 specifically.
	leftBoomStyleThick, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4F81BD"}, Pattern: 1},
		Font:      boomStyleBase.Font,
		Alignment: boomStyleBase.Alignment,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 5},
		},
	})
	f.SetCellStyle(sheetName, "D1", "E1", leftBoomStyleThick)

	// Update Protect function header (E2)
	protectColor := model.FunctionColors["Protect"]
	protectHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{protectColor}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 5},
		},
	})
	f.SetCellStyle(sheetName, "E2", "E2", protectHeaderStyle)

	// 5. Add Dependency Legend
	legendRow := currentRow + 1
	legendCell := fmt.Sprintf("A%d", legendRow)
	lastCol, _ := excelize.ColumnNumberToName(len(model.Functions) + 2)
	lastCell := fmt.Sprintf("%s%d", lastCol, legendRow)
	
	f.MergeCell(sheetName, legendCell, lastCell)
	f.SetCellValue(sheetName, legendCell, "Dependency emphasis: Govern is process-led | Identify/Protect lean Technology | Detect/Respond/Recover lean People")
	
	legendStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Italic: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheetName, legendCell, legendCell, legendStyle)

	// Auto-adjust column widths (approximated)
	f.SetColWidth(sheetName, "A", "B", 20)
	f.SetColWidth(sheetName, "C", "H", 28)

	return f.SaveAs(path)
}
