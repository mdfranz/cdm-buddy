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

	var assetClassesToAdd []string
	var editingInstance string // Track which instance we're editing (if any)

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
		// Add all asset classes for new assessments
		assetClassesToAdd = model.AssetClasses
	} else {
		fmt.Println(styleTitle.Render("\nResuming existing assessment..."))
		DisplaySummary(matrix)
		fmt.Println()

		var resumeChoice string
		resumeForm := huh.NewForm(huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
				Options(
					huh.NewOption("Add instances to a new asset class", "add"),
					huh.NewOption("Re-enter a specific instance", "edit"),
					huh.NewOption("Skip wizard and export", "export"),
				).
				Value(&resumeChoice),
		))
		if err := resumeForm.Run(); err != nil {
			return nil, err
		}

		switch resumeChoice {
		case "add":
			var assetChoice string
			assetOptions := make([]huh.Option[string], 0)
			for _, asset := range model.AssetClasses {
				assetOptions = append(assetOptions, huh.NewOption(asset, asset))
			}
			assetForm := huh.NewForm(huh.NewGroup(
				huh.NewSelect[string]().
					Title("Which asset class to add instances to?").
					Options(assetOptions...).
					Value(&assetChoice),
			))
			if err := assetForm.Run(); err != nil {
				return nil, err
			}
			assetClassesToAdd = []string{assetChoice}

		case "edit":
			var selectedInstance string
			instanceOptions := make([]huh.Option[string], 0)
			for _, asset := range model.AssetClasses {
				for _, inst := range matrix[asset] {
					label := fmt.Sprintf("%s / %s", asset, inst.Name)
					instanceOptions = append(instanceOptions, huh.NewOption(label, inst.Name))
				}
			}
			if len(instanceOptions) == 0 {
				fmt.Println("No instances to edit.")
				return matrix, nil
			}
			editForm := huh.NewForm(huh.NewGroup(
				huh.NewSelect[string]().
					Title("Which instance to re-enter?").
					Options(instanceOptions...).
					Value(&selectedInstance),
			))
			if err := editForm.Run(); err != nil {
				return nil, err
			}
			// Find the asset class for this instance
			for _, asset := range model.AssetClasses {
				for i, inst := range matrix[asset] {
					if inst.Name == selectedInstance {
						// Clear the instance and re-enter it
						matrix[asset][i] = model.AssetInstance{Name: selectedInstance, Cells: make(map[string]model.Cell)}
						assetClassesToAdd = []string{asset}
						editingInstance = selectedInstance
						break
					}
				}
			}

		case "export":
			// Skip wizard, go straight to export
			return matrix, nil
		}
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

	for _, asset := range assetClassesToAdd {
		for {
			var assetName string
			icon := model.AssetIcons[asset]
			assetHeader := styleHeader.Render(fmt.Sprintf("%s %s", icon, strings.ToUpper(asset)))
			assetDesc := styleDim.Render(model.AssetDescriptions[asset])

			// When editing, pre-fill the instance name
			if editingInstance != "" {
				assetName = editingInstance
			}

			// Skip the name prompt if we're editing
			if editingInstance == "" {
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

				assetName = strings.TrimSpace(assetName)
				if assetName == "" {
					break
				}

				if model.HasInstance(matrix, asset, assetName) {
					fmt.Printf("%s\n", styleRed.Render(fmt.Sprintf("❌ Instance '%s' already exists in %s. Choose a different name.", assetName, asset)))
					continue
				}
			} else {
				// When editing, just show a note and proceed
				fmt.Println(styleHeader.Render(fmt.Sprintf("Editing: %s / %s", asset, assetName)))
			}

			instance := model.AssetInstance{
				Name:  assetName,
				Cells: make(map[string]model.Cell),
			}

			// Offer to copy from existing instances in the same asset class (skip if editing)
			if editingInstance == "" {
				existingInstances := matrix[asset]
				if len(existingInstances) > 0 {
					var copyChoice string
					options := []huh.Option[string]{
						huh.NewOption("Start fresh", "fresh"),
					}
					for _, existing := range existingInstances {
						options = append(options, huh.NewOption("Copy from "+existing.Name, existing.Name))
					}

					copyForm := huh.NewForm(huh.NewGroup(
						huh.NewSelect[string]().
							Title(fmt.Sprintf("Prefill from existing %s instance?", asset)).
							Options(options...).
							Value(&copyChoice),
					))
					if err := copyForm.Run(); err == nil && copyChoice != "fresh" {
						// Find the instance to copy from
						for _, existing := range existingInstances {
							if existing.Name == copyChoice {
								// Deep copy the cells
								for funcName, cell := range existing.Cells {
									instance.Cells[funcName] = cell
								}
								break
							}
						}
					}
				}
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

				var markAsNone bool
				var cell model.Cell
				noneHint := styleDim.Render("(type 'none' for ❌)")

				confirmGroup := huh.NewGroup(
					huh.NewNote().Description(fmt.Sprintf("%s\n\n%s %s (%d/%d)\n%s\n%s %s", instanceHeader, funcHeader, boomTag, funcIdx+1, len(orderedFunctions), funcDesc, tipPrefix, tip)),
					huh.NewConfirm().
						Title("Mark all as ❌ (no coverage)?").
						Value(&markAsNone),
				)

				if err := huh.NewForm(confirmGroup).Run(); err != nil {
					if err == huh.ErrUserAborted {
						if quit := handleAbort(matrix); quit {
							return matrix, nil
						}
						continue
					}
					return nil, err
				}

				if markAsNone {
					instance.Cells[funcName] = model.Cell{Tech: "❌", People: "❌", Process: "❌"}
					continue
				}

				inputGroup := huh.NewGroup(
					huh.NewNote().Description(fmt.Sprintf("%s\n\n%s %s (%d/%d)\n%s\n%s %s", instanceHeader, funcHeader, boomTag, funcIdx+1, len(orderedFunctions), funcDesc, tipPrefix, tip)),
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

				inputForm := huh.NewForm(inputGroup)
				if err := inputForm.Run(); err != nil {
					if err == huh.ErrUserAborted {
						if quit := handleAbort(matrix); quit {
							return matrix, nil
						}
						continue
					}
					return nil, err
				}

				// Handle "none" -> ❌, and skip if all empty
				cell.Tech = formatValue(cell.Tech)
				cell.People = formatValue(cell.People)
				cell.Process = formatValue(cell.Process)

				instance.Cells[funcName] = cell
			}
			// Replace the instance if it already exists (for editing), otherwise append
			found := false
			for i, existing := range matrix[asset] {
				if existing.Name == instance.Name {
					matrix[asset][i] = instance
					found = true
					break
				}
			}
			if !found {
				matrix[asset] = append(matrix[asset], instance)
			}
			// After processing the first instance when editing, clear the flag
			if editingInstance != "" {
				editingInstance = ""
			}
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
	trimmed := strings.TrimSpace(s)
	if strings.ToLower(trimmed) == "none" {
		return "❌"
	}
	return strings.TrimSpace(s)
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
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).MarginBottom(1)
	fmt.Println(titleStyle.Render("\n 📊 Cyber Defense Matrix Coverage"))

	totalCells := 0
	filledCells := 0
	noneCells := 0

	// Collect all instances
	var allInstances []struct {
		asset    string
		instance string
	}
	for _, asset := range model.AssetClasses {
		for _, inst := range matrix[asset] {
			allInstances = append(allInstances, struct {
				asset    string
				instance string
			}{asset, inst.Name})
		}
	}

	if len(allInstances) == 0 {
		emptyStyle := lipgloss.NewStyle().Italic(true).Faint(true)
		fmt.Println(emptyStyle.Render("  (no instances added yet)"))
		return
	}

	// Build grid with box-drawing characters
	filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)   // Green ✓
	noneStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)      // Red ❌
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Faint(true)    // Gray -

	colWidth := 13
	instanceWidth := 32

	// Top border
	fmt.Print("  ┌" + strings.Repeat("─", instanceWidth-1))
	for range model.Functions {
		fmt.Print("┬" + strings.Repeat("─", colWidth-1))
	}
	fmt.Println("┐")

	// Header row
	fmt.Print("  │ " + padRight("Instance", instanceWidth-2))
	for _, funcName := range model.Functions {
		fmt.Print("│" + padCenter(funcName, colWidth-1))
	}
	fmt.Println("│")

	// Header separator
	fmt.Print("  ├" + strings.Repeat("─", instanceWidth-1))
	for range model.Functions {
		fmt.Print("┼" + strings.Repeat("─", colWidth-1))
	}
	fmt.Println("┤")

	// Data rows
	for i, item := range allInstances {
		fmt.Print("  │ " + padRight(item.instance, instanceWidth-2))
		instances := matrix[item.asset]
		var currentInst model.AssetInstance
		for _, inst := range instances {
			if inst.Name == item.instance {
				currentInst = inst
				break
			}
		}

		for _, funcName := range model.Functions {
			cell := currentInst.Cells[funcName]
			totalCells++

			var cellStatus, styledCell string
			if cell.Tech == "" && cell.People == "" && cell.Process == "" {
				cellStatus = "─"
				styledCell = emptyStyle.Render(cellStatus)
			} else if cell.Tech == "❌" && cell.People == "❌" && cell.Process == "❌" {
				cellStatus = "❌"
				styledCell = noneStyle.Render(cellStatus)
				noneCells++
			} else {
				cellStatus = "✓"
				styledCell = filledStyle.Render(cellStatus)
				filledCells++
			}
			fmt.Print("│" + padCenter(styledCell, colWidth-1))
		}
		fmt.Println("│")

		// Add separator between rows (not after last)
		if i < len(allInstances)-1 {
			fmt.Print("  ├" + strings.Repeat("─", instanceWidth-1))
			for range model.Functions {
				fmt.Print("┼" + strings.Repeat("─", colWidth-1))
			}
			fmt.Println("┤")
		}
	}

	// Bottom border
	fmt.Print("  └" + strings.Repeat("─", instanceWidth-1))
	for range model.Functions {
		fmt.Print("┴" + strings.Repeat("─", colWidth-1))
	}
	fmt.Println("┘")

	// Coverage summary
	coverage := 0
	if totalCells > 0 {
		coverage = (filledCells * 100) / totalCells
	}

	summaryStyle := lipgloss.NewStyle().
		MarginTop(1).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8"))

	summaryText := fmt.Sprintf(
		"  %s  %d/%d cells filled  •  %.0f%% coverage  •  %d marked no coverage",
		styleIcon.Render("📈"),
		filledCells, totalCells,
		float64(coverage),
		noneCells,
	)
	fmt.Println(summaryStyle.Render(summaryText))
}

func padRight(s string, width int) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Left).Render(s)
}

func padCenter(s string, width int) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(s)
}
