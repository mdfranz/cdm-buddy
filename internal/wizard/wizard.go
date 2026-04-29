package wizard

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cdmbuddy/internal/model"
	"cdmbuddy/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	styleTitle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Background(lipgloss.Color("24")).Padding(0, 2)
	styleSectionTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("110"))
	styleTag          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Background(lipgloss.Color("24")).Padding(0, 1)
	styleIcon         = lipgloss.NewStyle().MarginRight(1)
	styleDim          = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	styleYellow       = lipgloss.NewStyle().Foreground(lipgloss.Color("215"))
	styleBlue         = lipgloss.NewStyle().Foreground(lipgloss.Color("111"))
	styleRed          = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	styleGreen        = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
)

func getTerminalWidth() int {
	w, _, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil || w <= 0 {
		return 68
	}
	if w > 68 {
		return 68
	}
	return w
}

func getStyleCard() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(getTerminalWidth()).
		Padding(0, 0)
}

func renderText(text string) string {
	width := getTerminalWidth() - 2
	return lipgloss.NewStyle().
		Width(width).
		PaddingLeft(2).
		Foreground(lipgloss.Color("245")).
		Render(text)
}

func RunWizard(initial model.Matrix) (model.Matrix, bool, error) {
	// Clear screen for a clean start
	fmt.Print("\033[H\033[2J")

	matrix := initial
	if matrix == nil {
		matrix = model.EmptyMatrix()
	}

	// 1. Intro & Concepts
	isResuming := false
	for _, instances := range matrix {
		if len(instances) > 0 {
			isResuming = true
			break
		}
	}

	var assetClassesToAdd []string
	var editingInstance string

	if !isResuming {
		fmt.Print(renderIntroScreen())
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		// Add all asset classes for new assessments
		assetClassesToAdd = model.AssetClasses
	} else {
		fmt.Println(styleTitle.Render("\nResuming existing assessment..."))
		DisplaySummary(matrix)
		fmt.Println()

		var resumeChoice string
		resumeForm := ui.NewForm(huh.NewGroup(
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
			return nil, false, err
		}

		switch resumeChoice {
		case "add":
			var assetChoice string
			assetOptions := make([]huh.Option[string], 0)
			for _, asset := range model.AssetClasses {
				assetOptions = append(assetOptions, huh.NewOption(asset, asset))
			}
			assetForm := ui.NewForm(huh.NewGroup(
				huh.NewSelect[string]().
					Title("Which asset class to add instances to?").
					Options(assetOptions...).
					Value(&assetChoice),
			))
			if err := assetForm.Run(); err != nil {
				return nil, false, err
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
				return matrix, false, nil
			}
			editForm := ui.NewForm(huh.NewGroup(
				huh.NewSelect[string]().
					Title("Which instance to re-enter?").
					Options(instanceOptions...).
					Value(&selectedInstance),
			))
			if err := editForm.Run(); err != nil {
				return nil, false, err
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
			return matrix, false, nil
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

			// When editing, pre-fill the instance name
			if editingInstance != "" {
				assetName = editingInstance
			}

			// Skip the name prompt if we're editing
			if editingInstance == "" {
				nameForm := ui.NewForm(huh.NewGroup(
					huh.NewNote().Description(renderAssetCard(asset, icon)),
					huh.NewInput().
						Title(fmt.Sprintf("Add %s", model.AssetEntryLabels[asset])).
						Placeholder(model.AssetInstanceExamples[asset]+" (leave blank to finish)").
						Value(&assetName),
				))
				if err := nameForm.Run(); err != nil {
					if err == huh.ErrUserAborted {
						if quit := handleAbort(matrix); quit {
							return matrix, true, nil
						}
						continue
					}
					return nil, false, err
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
				fmt.Println(renderAssetCard(asset, icon) + "\n" + styleDim.Render("Editing: "+assetName))
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

					copyForm := ui.NewForm(huh.NewGroup(
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
					boomTag = "Cross-cutting"
					peopleInstruction = "e.g. Security Governance Lead, Risk Committee"
					processInstruction = "e.g. Risk Exception Process, Security Policy Review"
				} else {
					isRight := contains(model.RightOfBoomFunctions, funcName)
					if isRight {
						boomTag = "Right of Boom"
					} else {
						boomTag = "Left of Boom"
					}
					peopleInstruction = "e.g. SOC Analyst, SysAdmin"
					processInstruction = "e.g. Patch Management SOP, IR Plan"
				}

				techInstruction := model.TechExamples[asset+"-"+funcName]
				if techInstruction == "" {
					techInstruction = "e.g. vendor tool or platform name"
				}

				funcDesc := styleDim.Render(model.FunctionDescriptions[funcName])

				var tip string
				switch funcName {
				case "Govern":
					tip = "Think of this as the context layer: policy, oversight, and risk decisions around this asset."
				case "Identify", "Protect", "Recover":
					tip = "Map based on the primary asset being acted upon."
				case "Detect":
					tip = "Map based on the use case being monitored, not just the tool name."
				case "Respond":
					tip = "Map based on the asset being investigated or contained."
				}

				var markAsNone bool
				var cell model.Cell
				noneHint := styleDim.Render("Type 'none' in any field to mark explicit no coverage.")
				contextCard := renderFunctionCard(asset, assetName, icon, funcName, boomTag, funcIdx+1, len(orderedFunctions), funcDesc, tip)

				confirmGroup := huh.NewGroup(
					huh.NewNote().Description(contextCard),
					huh.NewConfirm().
						Title("No coverage for this function?").
						Value(&markAsNone),
				)

				if err := ui.NewForm(confirmGroup).Run(); err != nil {
					if err == huh.ErrUserAborted {
						if quit := handleAbort(matrix); quit {
							return matrix, true, nil
						}
						continue
					}
					return nil, false, err
				}

				if markAsNone {
					instance.Cells[funcName] = model.Cell{Tech: "❌", People: "❌", Process: "❌"}
					continue
				}

				inputGroup := huh.NewGroup(
					huh.NewNote().Description(contextCard),
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

				inputForm := ui.NewForm(inputGroup)
				if err := inputForm.Run(); err != nil {
					if err == huh.ErrUserAborted {
						if quit := handleAbort(matrix); quit {
							return matrix, true, nil
						}
						continue
					}
					return nil, false, err
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

	return matrix, false, nil
}

func handleAbort(matrix model.Matrix) bool {
	var quit bool
	confirmForm := ui.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Are you sure you want to quit?").
			Value(&quit),
	))
	_ = confirmForm.Run()
	if !quit {
		return false
	}

	saveChoice, err := ui.PromptSingleKeyChoice(
		"Save progress before exiting?",
		ui.SingleKeyOption{Key: 's', Label: "Save as JSON and Quit", Value: "json"},
		ui.SingleKeyOption{Key: 'q', Label: "Quit without saving", Value: "quit"},
	)
	if err != nil {
		fmt.Printf("Unable to read exit choice: %v\n", err)
		return false
	}

	if saveChoice == "json" {
		var path string
		pathForm := ui.NewForm(huh.NewGroup(
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
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("110")).MarginBottom(1)
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
	filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true) // Green ✓
	noneStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)    // Red ❌
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Faint(true)  // Gray -

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
		Padding(1, 2)

	summaryText := fmt.Sprintf(
		"  %s  %d/%d cells filled  •  %.0f%% coverage  •  %d marked no coverage",
		styleIcon.Render("📈"),
		filledCells, totalCells,
		float64(coverage),
		noneCells,
	)
	fmt.Println(summaryStyle.Render(summaryText))
}

func renderIntroScreen() string {
	lines := []string{
		styleTitle.Render("Cyber Defense Matrix (v1.1)"),
		getStyleCard().Render(lipgloss.JoinVertical(lipgloss.Left,
			styleSectionTitle.Render("Map coverage with intent"),
			renderText("Add concrete asset instances, then capture the primary Technology, People, and Process used across each function."),
			"",
			renderIntroStep("1", "Choose an asset class", "Services, Devices, Networks, Applications, Data, or Users."),
			renderIntroStep("2", "Add a real instance", "Use names like Customer Portal, Workstations, or Corp LAN."),
			renderIntroStep("3", "Map the primary control", "Each capability belongs in one cell. Avoid double-counting."),
			"",
			renderText(styleYellow.Render("Tip: ")+"type 'none' in any field to record explicit no coverage."),
		)),
		styleGreen.Render("Press Enter to start."),
	}
	return "\n" + strings.Join(lines, "\n") + "\n"
}

func renderIntroStep(number, title, detail string) string {
	stepTag := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Background(lipgloss.Color("236")).Padding(0, 1)
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center,
			stepTag.Render(number),
			" ",
			styleSectionTitle.Render(title),
		),
		renderText(detail),
		"",
	)
}

func renderAssetCard(asset, icon string) string {
	header := lipgloss.NewStyle().PaddingLeft(1).Render(
		lipgloss.JoinHorizontal(lipgloss.Center, styleTag.Render(icon+" "+asset), " ", styleSectionTitle.Render("Add an instance")),
	)
	return getStyleCard().Render(lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		renderText(model.AssetDescriptions[asset]),
		renderText("Examples: "+model.AssetInstanceExamples[asset]),
	))
}

func renderFunctionCard(asset, assetName, icon, functionName, stage string, step, total int, description, tip string) string {
	header1 := lipgloss.NewStyle().PaddingLeft(1).Render(
		lipgloss.JoinHorizontal(lipgloss.Center,
			styleTag.Render(icon+" "+asset),
			" ",
			styleSectionTitle.Render(assetName),
		),
	)
	header2 := lipgloss.NewStyle().PaddingLeft(1).Render(
		lipgloss.JoinHorizontal(lipgloss.Center,
			renderStagePill(functionName, lipgloss.Color("24")),
			" ",
			renderStageLabel(stage),
			" ",
			styleDim.Render(fmt.Sprintf("%d of %d", step, total)),
		),
	)

	lines := []string{
		header1,
		header2,
		"",
		renderText(description),
	}
	if tip != "" {
		lines = append(lines, "", renderText(styleYellow.Render("Tip: ")+tip))
	}
	return getStyleCard().Render(strings.Join(lines, "\n"))
}

func renderStageLabel(stage string) string {
	switch stage {
	case "Left of Boom":
		return renderStagePill(stage, lipgloss.Color("31"))
	case "Right of Boom":
		return renderStagePill(stage, lipgloss.Color("124"))
	default:
		return renderStagePill(stage, lipgloss.Color("94"))
	}
}

func renderStagePill(label string, background lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(background).
		Padding(0, 1).
		Render(label)
}

func padRight(s string, width int) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Left).Render(s)
}

func padCenter(s string, width int) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(s)
}
