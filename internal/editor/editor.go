package editor

import (
	"fmt"
	"os"
	"strings"

	"cdmbuddy/internal/model"
	"cdmbuddy/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	styleTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).MarginBottom(1)
	styleHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14")).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("14"))
	styleDim    = lipgloss.NewStyle().Italic(true).Faint(true)
	styleGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	styleRed    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

func RunEditor(matrix model.Matrix) (model.Matrix, error) {
	for {
		// 1. Main Menu: Asset Class Selection
		var selectedAssetClass string
		assetOptions := []huh.Option[string]{}
		for _, asset := range model.AssetClasses {
			count := len(matrix[asset])
			if count > 0 {
				label := fmt.Sprintf("%s (%d instances)", asset, count)
				assetOptions = append(assetOptions, huh.NewOption(label, asset))
			}
		}
		// If NO instances exist yet, allow selecting an asset class to add to
		if len(assetOptions) == 0 {
			for _, asset := range model.AssetClasses {
				assetOptions = append(assetOptions, huh.NewOption(asset, asset))
			}
		}

		assetOptions = append(assetOptions, huh.NewOption("---", "spacer"))
		assetOptions = append(assetOptions, huh.NewOption("Save & Exit", "exit"))

		mainMenu := ui.NewForm(huh.NewGroup(
			huh.NewSelect[string]().
				Title("Editor: Select Asset Class").
				Options(assetOptions...).
				Value(&selectedAssetClass).
				Filtering(true),
		))

		if err := mainMenu.Run(); err != nil {
			if err == huh.ErrUserAborted {
				if quit := handleAbort(matrix); quit {
					return matrix, nil
				}
				continue
			}
			return matrix, err
		}

		if selectedAssetClass == "spacer" {
			continue
		}

		if selectedAssetClass == "exit" {
			break
		}

		// 2. Instance Menu
		for {
			var selectedInstanceName string
			instanceOptions := []huh.Option[string]{}
			instanceOptions = append(instanceOptions, huh.NewOption("[Add New Instance...]", "add_new"))
			instanceOptions = append(instanceOptions, huh.NewOption("---", "spacer"))

			instances := matrix[selectedAssetClass]
			for _, inst := range instances {
				instanceOptions = append(instanceOptions, huh.NewOption(inst.Name, inst.Name))
			}
			instanceOptions = append(instanceOptions, huh.NewOption("---", "spacer"))
			instanceOptions = append(instanceOptions, huh.NewOption("Back to Asset Classes", "back"))

			instanceMenu := ui.NewForm(huh.NewGroup(
				huh.NewSelect[string]().
					Title(fmt.Sprintf("%s: Select Instance", selectedAssetClass)).
					Options(instanceOptions...).
					Value(&selectedInstanceName).
					Filtering(true),
			))

			if err := instanceMenu.Run(); err != nil {
				if err == huh.ErrUserAborted {
					if quit := handleAbort(matrix); quit {
						return matrix, nil
					}
					continue
				}
				return matrix, err
			}

			if selectedInstanceName == "spacer" {
				continue
			}

			if selectedInstanceName == "back" {
				break
			}

			if selectedInstanceName == "add_new" {
				var newName string
				nameForm := ui.NewForm(huh.NewGroup(
					huh.NewInput().
						Title(fmt.Sprintf("Name for new %s instance", selectedAssetClass)).
						Placeholder(model.AssetInstanceExamples[selectedAssetClass]).
						Value(&newName),
				))
				if err := nameForm.Run(); err != nil {
					if err == huh.ErrUserAborted {
						if quit := handleAbort(matrix); quit {
							return matrix, nil
						}
					}
					continue
				}

				newName = strings.TrimSpace(newName)
				if newName == "" {
					continue
				}

				if model.HasInstance(matrix, selectedAssetClass, newName) {
					fmt.Printf("\n%s\n", styleRed.Render("❌ That instance already exists."))
					continue
				}

				// Create it
				newInstance := model.AssetInstance{
					Name:  newName,
					Cells: make(map[string]model.Cell),
				}
				matrix[selectedAssetClass] = append(matrix[selectedAssetClass], newInstance)
				selectedInstanceName = newName
				fmt.Printf("\n%s\n", styleGreen.Render(fmt.Sprintf("✓ Created new instance: %s", newName)))
			}

			// Find the actual instance pointer/index
			var instIdx int
			found := false
			for i, inst := range matrix[selectedAssetClass] {
				if inst.Name == selectedInstanceName {
					instIdx = i
					found = true
					break
				}
			}

			if !found {
				continue
			}

			// 3. Function/Cell Menu
			for {
				var selectedFunction string
				functionOptions := []huh.Option[string]{}
				for _, fn := range model.Functions {
					cell := matrix[selectedAssetClass][instIdx].Cells[fn]
					preview := "empty"
					if cell.Tech != "" || cell.People != "" || cell.Process != "" {
						preview = fmt.Sprintf("[%s, %s, %s]", truncate(cell.Tech, 15), truncate(cell.People, 15), truncate(cell.Process, 15))
					}
					label := fmt.Sprintf("%s: %s", fn, styleDim.Render(preview))
					functionOptions = append(functionOptions, huh.NewOption(label, fn))
				}
				functionOptions = append(functionOptions, huh.NewOption("---", "spacer"))
				functionOptions = append(functionOptions, huh.NewOption("Delete This Instance", "delete"))
				functionOptions = append(functionOptions, huh.NewOption("Back to Instances", "back"))

				functionMenu := ui.NewForm(huh.NewGroup(
					huh.NewSelect[string]().
						Title(fmt.Sprintf("Editing %s / %s: Select Function", selectedAssetClass, selectedInstanceName)).
						Options(functionOptions...).
						Value(&selectedFunction),
				))

				if err := functionMenu.Run(); err != nil {
					if err == huh.ErrUserAborted {
						if quit := handleAbort(matrix); quit {
							return matrix, nil
						}
						continue
					}
					return matrix, err
				}

				if selectedFunction == "spacer" {
					continue
				}

				if selectedFunction == "back" {
					break
				}

				if selectedFunction == "delete" {
					var confirmDelete bool
					confirmDeleteForm := ui.NewForm(huh.NewGroup(
						huh.NewConfirm().
							Title(fmt.Sprintf("Are you sure you want to delete '%s' and ALL its data?", selectedInstanceName)).
							Value(&confirmDelete),
					))
					if err := confirmDeleteForm.Run(); err != nil {
						continue
					}

					if confirmDelete {
						// Filter out the instance
						instances := matrix[selectedAssetClass]
						newInstances := []model.AssetInstance{}
						for i, inst := range instances {
							if i != instIdx {
								newInstances = append(newInstances, inst)
							}
						}
						matrix[selectedAssetClass] = newInstances
						fmt.Printf("\n%s\n", styleRed.Render(fmt.Sprintf("🗑️ Deleted instance: %s", selectedInstanceName)))
						break // Exit Function Menu
					}
					continue
				}

				// 4. Edit Cell Form
				currentCell := matrix[selectedAssetClass][instIdx].Cells[selectedFunction]
				newCell := currentCell

				var action string = "save"
				editForm := ui.NewForm(
					huh.NewGroup(
						huh.NewNote().Title(fmt.Sprintf("Editing %s / %s / %s", selectedAssetClass, selectedInstanceName, selectedFunction)),
						huh.NewInput().
							Title("Technology/Vendor").
							Value(&newCell.Tech),
						huh.NewInput().
							Title("People/Responsible Role").
							Value(&newCell.People),
						huh.NewInput().
							Title("Process/Procedure").
							Value(&newCell.Process),
					),
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("Action").
							Options(
								huh.NewOption("Save Changes", "save"),
								huh.NewOption("Clear Cell Content", "clear"),
								huh.NewOption("Discard Changes", "discard"),
							).
							Value(&action),
					),
				)

				if err := editForm.Run(); err != nil {
					if err == huh.ErrUserAborted {
						if quit := handleAbort(matrix); quit {
							return matrix, nil
						}
						continue
					}
					return matrix, err
				}

				switch action {
				case "clear":
					matrix[selectedAssetClass][instIdx].Cells[selectedFunction] = model.Cell{}
					fmt.Println(styleRed.Render(fmt.Sprintf("\n🗑️ Cleared cell: %s", selectedFunction)))
				case "save":
					// Update the matrix with new values
					matrix[selectedAssetClass][instIdx].Cells[selectedFunction] = newCell
					fmt.Println(styleGreen.Render(fmt.Sprintf("\n✓ Saved %s / %s / %s", selectedAssetClass, selectedInstanceName, selectedFunction)))
				case "discard":
					fmt.Println(styleDim.Render("\n(Changes discarded)"))
				}
			}
		}
	}

	return matrix, nil
}

func handleAbort(matrix model.Matrix) bool {
	var quit bool
	confirmForm := ui.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Are you sure you want to quit the editor?").
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
		os.Exit(0)
	} else if saveChoice == "quit" {
		os.Exit(0)
	}
	return true
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
