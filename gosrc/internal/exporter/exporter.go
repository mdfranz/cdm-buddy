package exporter

import (
	"fmt"

	"cdmbuddy/internal/model"
	"github.com/xuri/excelize/v2"
)

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
	f.SetCellValue(sheetName, "B1", "GOVERN (Cross-Cutting)")
	f.SetCellStyle(sheetName, "B1", "B1", governBoomStyle)

	f.MergeCell(sheetName, "C1", "D1")
	f.SetCellValue(sheetName, "C1", "LEFT OF BOOM (Pre-Event)")
	f.SetCellStyle(sheetName, "C1", "D1", leftBoomStyle)

	f.MergeCell(sheetName, "E1", "G1")
	f.SetCellValue(sheetName, "E1", "RIGHT OF BOOM (Post-Event)")
	f.SetCellStyle(sheetName, "E1", "G1", rightBoomStyle)

	// 2. Add Function Headers
	for i, funcName := range model.Functions {
		col, _ := excelize.ColumnNumberToName(i + 2)
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

	// 3. Add Asset Headers and Data
	for r, asset := range model.AssetClasses {
		rowIdx := r + 3
		cellName := fmt.Sprintf("A%d", rowIdx)
		f.SetCellValue(sheetName, cellName, asset)
		f.SetCellStyle(sheetName, cellName, cellName, headerStyle)

		for c, funcName := range model.Functions {
			col, _ := excelize.ColumnNumberToName(c + 2)
			cell := fmt.Sprintf("%s%d", col, rowIdx)
			
			cellData := matrix[asset][funcName]
			value := fmt.Sprintf("TECH: %s\nPEOPLE: %s\nPROCESS: %s", cellData.Tech, cellData.People, cellData.Process)
			f.SetCellValue(sheetName, cell, value)

			color := model.FunctionColors[funcName]
			
			// Check if we need thick right border for "Protect" (column D)
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
	}

	// 4. Thick border for Boom Boundary
	// Protect is column D (4). We apply it to rows 1 and 2 specifically.
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
	f.SetCellStyle(sheetName, "C1", "D1", leftBoomStyleThick)

	// Update Protect function header (D2)
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
	f.SetCellStyle(sheetName, "D2", "D2", protectHeaderStyle)

	// 5. Add Dependency Legend
	legendRow := len(model.AssetClasses) + 4
	legendCell := fmt.Sprintf("A%d", legendRow)
	lastCol, _ := excelize.ColumnNumberToName(len(model.Functions) + 1)
	lastCell := fmt.Sprintf("%s%d", lastCol, legendRow)
	
	f.MergeCell(sheetName, legendCell, lastCell)
	f.SetCellValue(sheetName, legendCell, "Dependency emphasis: Govern is process-led | Identify/Protect lean Technology | Detect/Respond/Recover lean People")
	
	legendStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Italic: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheetName, legendCell, legendCell, legendStyle)

	// Auto-adjust column widths (approximated)
	f.SetColWidth(sheetName, "A", "A", 15)
	f.SetColWidth(sheetName, "B", "G", 28)

	return f.SaveAs(path)
}
