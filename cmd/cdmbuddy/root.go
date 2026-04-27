package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"cdmbuddy/internal/editor"
	"cdmbuddy/internal/exporter"
	"cdmbuddy/internal/model"
	"cdmbuddy/internal/wizard"
	"github.com/spf13/cobra"
)

var (
	inputPath  string
	editPath   string
	outputPath string
	csvPath    string
	jsonPath   string
	assessName string
	reportOnly bool
)

var rootCmd = &cobra.Command{
	Use:   "cdmbuddy",
	Short: "Cyber Defense Matrix Wizard",
	Long:  `A guided tool for building your Cyber Defense Matrix, inspired by Sounil Yu's framework.`,
	Run: func(cmd *cobra.Command, args []string) {
		var matrix model.Matrix
		var err error

		// Handle Edit mode separately
		if editPath != "" {
			fmt.Printf("Loading data for editing from: %s\n", editPath)
			matrix, err = model.LoadFromJson(editPath)
			if err != nil {
				fmt.Printf("Error loading JSON: %v\n", err)
				os.Exit(1)
			}

			// In edit mode, we use the input filename to determine the base for other exports
			if assessName == "" {
				base := strings.TrimSuffix(editPath, ".json")
				base = strings.TrimSuffix(base, "_cdm")
				assessName = base
			}

			// Run the dedicated Editor instead of the Wizard
			matrix, err = editor.RunEditor(matrix)
			if err != nil {
				fmt.Printf("Error running editor: %v\n", err)
				os.Exit(1)
			}
		} else if inputPath != "" {
			fmt.Printf("Loading data from: %s\n", inputPath)
			matrix, err = model.LoadFromJson(inputPath)
			if err != nil {
				fmt.Printf("Error loading JSON: %v\n", err)
				os.Exit(1)
			}
		}

		if editPath == "" && !reportOnly {
			// Always run the wizard to allow adding more data/reviewing unless --report is set
			if inputPath == "" {
				// If name not provided via flag, prompt for it
				if assessName == "" {
					fmt.Print("Enter assessment name: ")
					reader := bufio.NewReader(os.Stdin)
					input, _ := reader.ReadString('\n')
					assessName = strings.TrimSpace(input)
					if assessName == "" {
						assessName = fmt.Sprintf("CDM_%s", time.Now().Format("060102150405"))
					}
				}
				assessName = strings.ReplaceAll(assessName, " ", "_")
			}
		}

		// Shared path logic
		if assessName != "" {
			if jsonPath == "cdm_data.json" {
				jsonPath = fmt.Sprintf("%s_cdm.json", assessName)
			}
			if csvPath == "cdm_output.csv" {
				csvPath = fmt.Sprintf("%s_cdm.csv", assessName)
			}
			if outputPath == "cdm_output.xlsx" {
				outputPath = fmt.Sprintf("%s_cdm.xlsx", assessName)
			}
		}

		if editPath == "" && !reportOnly {
			var aborted bool
			matrix, aborted, err = wizard.RunWizard(matrix)
			if err != nil {
				fmt.Printf("Error running wizard: %v\n", err)
				os.Exit(1)
			}
			if aborted {
				return
			}
		}

		wizard.DisplaySummary(matrix)

		err = exporter.ExportToExcel(matrix, outputPath)
		if err != nil {
			fmt.Printf("Error exporting to Excel: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nSuccess! Your Cyber Defense Matrix has been exported to %s\n", outputPath)

		if csvPath != "" {
			err = exporter.ExportToCSV(matrix, csvPath)
			if err != nil {
				fmt.Printf("Error exporting to CSV: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("CSV export: Data saved to %s\n", csvPath)
		}

		if jsonPath != "" {
			err = model.SaveToJson(matrix, jsonPath)
			if err != nil {
				fmt.Printf("Error exporting to JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("JSON export: Data saved to %s\n", jsonPath)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Load CDM data from a JSON file")
	rootCmd.Flags().StringVarP(&editPath, "edit", "e", "", "Browse and edit an existing JSON file")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "cdm_output.xlsx", "Output Excel file path")
	rootCmd.Flags().StringVarP(&csvPath, "csv", "c", "cdm_output.csv", "Output CSV file path")
	rootCmd.Flags().StringVar(&jsonPath, "export-json", "cdm_data.json", "Also export CDM data as JSON")
	rootCmd.Flags().StringVarP(&assessName, "name", "n", "", "Assessment name (used for filename)")
	rootCmd.Flags().BoolVarP(&reportOnly, "report", "r", false, "Generate report and exit without showing the wizard (requires -i)")
}
