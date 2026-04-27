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

func RunWizard(initial model.Matrix) (model.Matrix, error) {
	matrix := initial
	if matrix == nil {
		matrix = model.EmptyMatrix()
	}

	// 1. Intro & Concepts (Skip if resuming)
	isResuming := false
	for _, instances := range matrix {
		if len(instances) > 0 {
			isResuming = true
			break
		}
	}

	if !isResuming {
		mission := huh.NewNote().
			Title("The Mission: Map Your Defense").
			Description("The **Cyber Defense Matrix** reveals where your security is redundant and where it is completely absent.\n\n" +
				"**1. Categories (Asset Classes)**\n" +
				"Think of these as the 'folders' or broad rows in the matrix: Services, Devices, Networks, Apps, Data, and Users.\n\n" +
				"**2. Specific Items (Asset Instances)**\n" +
				"Within each category, you add the specific things you care about. For example:\n" +
				"• Under **Devices**, you might add 'Workstations' and 'Servers'.\n" +
				"• Under **Services**, you might add 'Salesforce' and 'Atlassian'.\n\n" +
				"**3. Mapping the Defense**\n" +
				"For every specific item, you will enter the **Technology, People, and Process** used to defend it across 6 functions (Identify, Protect, etc.).")

		axisNote := huh.NewNote().
			Title("Rules of the Matrix").
			Description(
				"**The MECE Rule**: Every capability belongs in **exactly one cell**. Map to the **primary** function only.\n\n" +
					"*(type 'none' in any field to explicitly mark a ❌ gap)*")

		introForm := huh.NewForm(
			huh.NewGroup(mission),
			huh.NewGroup(axisNote),
		)

		if err := introForm.Run(); err != nil {
			return nil, err
		}
	} else {
		fmt.Println(styleTitle.Render("\nResuming existing assessment..."))
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
					if quit := handleAbort(matrix); quit {
						return matrix, nil
					}
					continue
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

				var status string = "details"
				statusGroup := huh.NewGroup(
					huh.NewNote().Description(fmt.Sprintf("%s\n\n%s %s (%d/%d)\n%s\n%s %s", instanceHeader, funcHeader, boomTag, funcIdx+1, len(orderedFunctions), funcDesc, tipPrefix, tip)),
					huh.NewSelect[string]().
						Title(fmt.Sprintf("Configure %s?", funcName)).
						Options(
							huh.NewOption("Enter details", "details"),
							huh.NewOption("Mark all as 'None' (❌)", "none"),
							huh.NewOption("Skip this function", "skip"),
						).
						Value(&status),
				)

				if err := huh.NewForm(statusGroup).Run(); err != nil {
					if err == huh.ErrUserAborted {
						if quit := handleAbort(matrix); quit {
							return matrix, nil
						}
						// If they didn't quit, we re-prompt for the same status group
						// A bit tricky in this nested loop structure, but standard break/continue won't work perfectly.
						// Simplest is to just allow them to 'skip' the cell if they cancel the cell-level abort.
						continue
					}
					return nil, err
				}

				if status == "none" {
					instance.Cells[funcName] = model.Cell{Tech: "❌", People: "❌", Process: "❌"}
					continue
				}
				if status == "skip" {
					instance.Cells[funcName] = model.Cell{}
					continue
				}

				cell := model.Cell{}
				noneHint := styleDim.Render("(type 'none' for ❌)")

				group := huh.NewGroup(
					huh.NewInput().
						Title("Technology/Vendor").
						Placeholder(techInstruction).
						Suggestions([]string{"none"}).
						Value(&cell.Tech),
					huh.NewInput().
						Title("People/Responsible Role").
						Placeholder(peopleInstruction).
						Suggestions([]string{"none"}).
						Value(&cell.People),
					huh.NewInput().
						Title("Process/Procedure").
						Placeholder(processInstruction).
						Suggestions([]string{"none"}).
						Value(&cell.Process),
					huh.NewNote().Description(noneHint),
				)

				form := huh.NewForm(group)
				if err := form.Run(); err != nil {
					if err == huh.ErrUserAborted {
						if quit := handleAbort(matrix); quit {
							return matrix, nil
						}
						continue
					}
					return nil, err
				}

				// Handle "none" -> ❌
				cell.Tech = formatValue(cell.Tech)
				cell.People = formatValue(cell.People)
				cell.Process = formatValue(cell.Process)

				instance.Cells[funcName] = cell
			}
			matrix[asset] = append(matrix[asset], instance)
		}
	}

	return matrix, nil
}

func handleAbort(matrix model.Matrix) bool {
	var quit bool
	confirmForm := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Are you sure you want to quit?").
			Value(&quit),
	))
	_ = confirmForm.Run()
	if !quit {
		return false
	}

	var saveChoice string
	saveForm := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Save progress before exiting?").
			Options(
				huh.NewOption("Save as JSON and Quit", "json"),
				huh.NewOption("Quit without saving", "quit"),
			).
			Value(&saveChoice),
	))
	_ = saveForm.Run()

	if saveChoice == "json" {
		var path string
		pathForm := huh.NewForm(huh.NewGroup(
			huh.NewInput().Title("JSON output path").Value(&path).Placeholder("cdm_progress.json"),
		))
		_ = pathForm.Run()
		if path == "" {
			path = "cdm_progress.json"
		}
		_ = model.SaveToJson(matrix, path)
		fmt.Printf("Progress saved to %s\n", path)
	}
	return true
}

func formatValue(s string) string {
	if strings.ToLower(strings.TrimSpace(s)) == "none" {
		return "❌"
	}
	return s
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
