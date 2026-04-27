package main

import (
	"fmt"
	"os"

	"cdmbuddy/internal/exporter"
	"cdmbuddy/internal/model"
	"cdmbuddy/internal/wizard"
	"github.com/spf13/cobra"
)

var (
	inputPath  string
	outputPath string
	jsonPath   string
)

var rootCmd = &cobra.Command{
	Use:   "cdmbuddy",
	Short: "Cyber Defense Matrix Wizard",
	Long:  `A guided tool for building your Cyber Defense Matrix, inspired by Sounil Yu's framework.`,
	Run: func(cmd *cobra.Command, args []string) {
		var matrix model.Matrix
		var err error

		if inputPath != "" {
			fmt.Printf("Loading data from: %s\n", inputPath)
			matrix, err = model.LoadFromJson(inputPath)
			if err != nil {
				fmt.Printf("Error loading JSON: %v\n", err)
				os.Exit(1)
			}
		} else {
			matrix, err = wizard.RunWizard()
			if err != nil {
				fmt.Printf("Error running wizard: %v\n", err)
				os.Exit(1)
			}
		}

		wizard.DisplaySummary(matrix)

		err = exporter.ExportToExcel(matrix, outputPath)
		if err != nil {
			fmt.Printf("Error exporting to Excel: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nSuccess! Your Cyber Defense Matrix has been exported to %s\n", outputPath)

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
	rootCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Load CDM data from a JSON file instead of running the interactive wizard")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "cdm_output.xlsx", "Output Excel file path")
	rootCmd.Flags().StringVar(&jsonPath, "export-json", "", "Also export CDM data as JSON (useful for saving and resuming later)")
}
