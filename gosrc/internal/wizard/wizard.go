package wizard

import (
	"fmt"
	"os"
	"strings"

	"cdmbuddy/internal/model"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	styleTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).MarginBottom(1)
	styleHeader  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14")).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("14"))
	styleIcon    = lipgloss.NewStyle().MarginRight(1)
	styleDim     = lipgloss.NewStyle().Italic(true).Faint(true)
	styleYellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	styleBlue    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	styleRed     = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	styleReverse = lipgloss.NewStyle().Reverse(true).Bold(true).Padding(0, 1)
)

func RunWizard() (model.Matrix, error) {
	matrix := model.EmptyMatrix()

	// 1. Intro
	intro := huh.NewNote().
		Title("Cyber Defense Matrix Wizard").
		Description("Guided by Sounil Yu's 'The Essential Guide to Navigating the Cybersecurity Landscape'\n\nThe CDM is a MECE framework designed to provide structural clarity and identify gaps in your security posture.")

	concepts := huh.NewNote().
		Title("Core Concepts").
		Description(
			"• Left of Boom (Proactive): Identify and Protect. Focuses on Structural Awareness.\n" +
				"• Right of Boom (Reactive): Detect, Respond, and Recover. Focuses on Situational Awareness.\n\n" +
				"The Dependency Continuum:\n" +
				"• Technology dependency is strongest in Identify/Protect.\n" +
				"• People dependency grows significantly in Detect/Respond/Recover.\n" +
				"• Process is the essential, consistent foundation throughout all functions.",
		)

	form := huh.NewForm(huh.NewGroup(intro, concepts))
	if err := form.Run(); err != nil {
		return nil, err
	}

	// 2. Iterative Wizard
	orderedFunctions := []string{}
	for _, f := range model.Functions {
		if !contains(model.GovernFunctions, f) {
			orderedFunctions = append(orderedFunctions, f)
		}
	}
	for _, f := range model.Functions {
		if contains(model.GovernFunctions, f) {
			orderedFunctions = append(orderedFunctions, f)
		}
	}

	totalCells := len(model.AssetClasses) * len(orderedFunctions)
	cellNum := 0

	for assetIdx, asset := range model.AssetClasses {
		for funcIdx, funcName := range orderedFunctions {
			cellNum++
			progress := fmt.Sprintf("(Asset %d/%d | Function %d/%d | Cell %d/%d)", assetIdx+1, len(model.AssetClasses), funcIdx+1, len(orderedFunctions), cellNum, totalCells)

			var boomTag string
			var peopleInstruction string
			var processInstruction string

			if contains(model.GovernFunctions, funcName) {
				boomTag = styleYellow.Render("Cross-cutting")
				peopleInstruction = "e.g. Security Governance Lead, Risk Committee"
				processInstruction = "e.g. Risk Exception Process, Security Policy Review"
			} else {
				isRight := contains(model.RightOfBoomFunctions, funcName)
				if isRight {
					boomTag = styleRed.Render("Right of Boom")
				} else {
					boomTag = styleBlue.Render("Left of Boom")
				}
				peopleInstruction = "e.g. SOC Analyst, SysAdmin"
				processInstruction = "e.g. Patch Management SOP, IR Plan"
			}

			techInstruction := model.TechExamples[asset+"-"+funcName]
			if techInstruction == "" {
				techInstruction = "e.g. vendor tool or platform name"
			}

			icon := model.AssetIcons[asset]
			assetHeader := styleHeader.Render(fmt.Sprintf("%s %s", icon, strings.ToUpper(asset)))
			assetDesc := styleDim.Render(model.AssetDescriptions[asset])
			funcHeader := styleReverse.Render(" " + funcName + " ")
			funcDesc := styleDim.Render(model.FunctionDescriptions[funcName])

			var tip string
			tipPrefix := styleYellow.Render("💡 Tip:")
			switch funcName {
			case "Govern":
				tip = styleDim.Render("Think of this as the 'Context'—policies, risk appetite, and oversight governing this asset class.")
			case "Identify", "Protect", "Recover":
				tip = styleDim.Render("Map based on the primary asset being acted upon.")
			case "Detect":
				tip = styleDim.Render("Map based on the Use Case (e.g., Insider Threat maps to Users).")
			case "Respond":
				tip = styleDim.Render("Map based on the asset being investigated or contained.")
			}

			cell := matrix[asset][funcName]
			
			// Creating a group for this specific cell
			group := huh.NewGroup(
				huh.NewNote().Description(fmt.Sprintf("%s\n%s\n\n%s %s %s\n%s\n%s %s", assetHeader, assetDesc, funcHeader, boomTag, styleDim.Render(progress), funcDesc, tipPrefix, tip)),
				huh.NewInput().
					Title("Technology/Vendor").
					Placeholder(techInstruction).
					Value(&cell.Tech),
				huh.NewInput().
					Title("People/Responsible Role").
					Placeholder(peopleInstruction).
					Value(&cell.People),
				huh.NewInput().
					Title("Process/Procedure").
					Placeholder(processInstruction).
					Value(&cell.Process),
			)

			form := huh.NewForm(group)
			err := form.Run()
			
			if err != nil {
				// Handle interrupt
				if err == huh.ErrUserAborted {
					var quit bool
					confirmForm := huh.NewForm(huh.NewGroup(
						huh.NewConfirm().
							Title("Are you sure you want to quit?").
							Value(&quit),
					))
					_ = confirmForm.Run()
					if quit {
						var saveChoice string
						saveForm := huh.NewForm(huh.NewGroup(
							huh.NewSelect[string]().
								Title("Save progress before exiting?").
								Options(
									huh.NewOption("Save as Excel", "excel"),
									huh.NewOption("Save as JSON", "json"),
									huh.NewOption("Quit without saving", "quit"),
								).
								Value(&saveChoice),
						))
						_ = saveForm.Run()
						if saveChoice == "excel" {
							return matrix, nil
						} else if saveChoice == "json" {
							var path string
							pathForm := huh.NewForm(huh.NewGroup(
								huh.NewInput().Title("JSON output path").Value(&path).Placeholder("cdm_progress.json"),
							))
							_ = pathForm.Run()
							if path == "" { path = "cdm_progress.json" }
							_ = model.SaveToJson(matrix, path)
							fmt.Printf("Progress saved to %s\n", path)
							os.Exit(0)
						} else {
							fmt.Println("Quitting without saving.")
							os.Exit(0)
						}
					}
					// If they didn't quit, we loop back to the same cell by decrementing or just letting it continue
					// But huh doesn't easily allow jumping back in a simple loop without re-running.
					// Let's just continue the loop which will effectively re-ask this cell since cell data is passed by pointer.
					// Actually we need to make sure 'cell' is updated in the matrix.
					cellNum-- 
					continue 
				}
				return nil, err
			}
			matrix[asset][funcName] = cell
		}
	}

	return matrix, nil
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func DisplaySummary(matrix model.Matrix) {
	// Simple summary for now, we can enhance with a table library later if needed
	fmt.Println("\nCyber Defense Matrix Summary")
	fmt.Println(strings.Repeat("-", 80))
	
	for _, asset := range model.AssetClasses {
		fmt.Printf("Asset: %s\n", asset)
		for _, funcName := range model.Functions {
			cell := matrix[asset][funcName]
			if cell.Tech != "" || cell.People != "" || cell.Process != "" {
				fmt.Printf("  [%s] Tech: %s, People: %s, Process: %s\n", funcName, cell.Tech, cell.People, cell.Process)
			}
		}
	}
	
	// Gap Analysis
	fullCount := 0
	partialCount := 0
	emptyCount := 0
	var emptyCells []string
	var partialCells []string

	for _, asset := range model.AssetClasses {
		for _, funcName := range model.Functions {
			cell := matrix[asset][funcName]
			filled := 0
			if cell.Tech != "" { filled++ }
			if cell.People != "" { filled++ }
			if cell.Process != "" { filled++ }

			if filled == 3 {
				fullCount++
			} else if filled == 0 {
				emptyCount++
				emptyCells = append(emptyCells, fmt.Sprintf("%s/%s", asset, funcName))
			} else {
				partialCount++
				partialCells = append(partialCells, fmt.Sprintf("%s/%s (%d/3 fields)", asset, funcName, filled))
			}
		}
	}

	totalCells := len(model.AssetClasses) * len(model.Functions)
	fmt.Printf("\nGap Analysis: %d/%d cells fully populated. %d partial, %d empty.\n", fullCount, totalCells, partialCount, emptyCount)
	
	if len(emptyCells) > 0 {
		fmt.Println("Empty cells:")
		for _, c := range emptyCells {
			fmt.Printf("  - %s\n", c)
		}
	}
}
