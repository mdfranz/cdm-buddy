package wizard

import (
	"fmt"
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

	for _, asset := range model.AssetClasses {
		for {
			var assetName string
			icon := model.AssetIcons[asset]
			assetHeader := styleHeader.Render(fmt.Sprintf("%s %s", icon, strings.ToUpper(asset)))
			assetDesc := styleDim.Render(model.AssetDescriptions[asset])

			nameForm := huh.NewForm(huh.NewGroup(
				huh.NewNote().Description(fmt.Sprintf("%s\n%s", assetHeader, assetDesc)),
				huh.NewInput().
					Title(fmt.Sprintf("Add a specific %s instance?", asset)).
					Placeholder("e.g. Workstations, Servers (leave blank to skip/finish)").
					Value(&assetName),
			))
			if err := nameForm.Run(); err != nil {
				if err == huh.ErrUserAborted {
					return matrix, nil
				}
				return nil, err
			}

			if assetName == "" {
				break
			}

			instance := model.AssetInstance{
				Name:  assetName,
				Cells: make(map[string]model.Cell),
			}

			for funcIdx, funcName := range orderedFunctions {
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

				instanceHeader := styleHeader.Render(fmt.Sprintf("%s %s: %s", icon, strings.ToUpper(asset), assetName))
				funcHeader := styleReverse.Render(" " + funcName + " ")
				funcDesc := styleDim.Render(model.FunctionDescriptions[funcName])

				var tip string
				tipPrefix := styleYellow.Render("💡 Tip:")
				switch funcName {
				case "Govern":
					tip = styleDim.Render("Think of this as the 'Context'—policies, risk appetite, and oversight governing this asset.")
				case "Identify", "Protect", "Recover":
					tip = styleDim.Render("Map based on the primary asset being acted upon.")
				case "Detect":
					tip = styleDim.Render("Map based on the Use Case (e.g., Insider Threat maps to Users).")
				case "Respond":
					tip = styleDim.Render("Map based on the asset being investigated or contained.")
				}

				cell := model.Cell{}

				group := huh.NewGroup(
					huh.NewNote().Description(fmt.Sprintf("%s\n\n%s %s (%d/%d)\n%s\n%s %s", instanceHeader, funcHeader, boomTag, funcIdx+1, len(orderedFunctions), funcDesc, tipPrefix, tip)),
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
				if err := form.Run(); err != nil {
					if err == huh.ErrUserAborted {
						// On abort during instance entry, we just finish this instance and ask if they want to save or continue
						// For simplicity, let's just break and return current matrix
						return matrix, nil
					}
					return nil, err
				}
				instance.Cells[funcName] = cell
			}
			matrix[asset] = append(matrix[asset], instance)
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
	fmt.Println("\nCyber Defense Matrix Summary")
	fmt.Println(strings.Repeat("-", 80))

	totalInstances := 0
	populatedCells := 0

	for _, asset := range model.AssetClasses {
		instances := matrix[asset]
		if len(instances) == 0 {
			continue
		}
		fmt.Printf("Asset Class: %s\n", asset)
		for _, instance := range instances {
			totalInstances++
			fmt.Printf("  Instance: %s\n", instance.Name)
			for _, funcName := range model.Functions {
				cell := instance.Cells[funcName]
				if cell.Tech != "" || cell.People != "" || cell.Process != "" {
					populatedCells++
					fmt.Printf("    [%s] Tech: %s, People: %s, Process: %s\n", funcName, cell.Tech, cell.People, cell.Process)
				}
			}
		}
	}

	fmt.Printf("\nSummary: %d asset instances defined. %d total populated cells.\n", totalInstances, populatedCells)
}
