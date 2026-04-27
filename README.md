![CDM Buddy Logo](cdm-buddy.png)

# Cyber Defense Matrix (CDM) Buddy

An interactive CLI tool to map your security portfolio to Sounil Yu's Cyber Defense Matrix.

CDM Buddy helps security teams visualize their coverage, identify gaps, and understand the balance between technology, people, and process across the NIST Cybersecurity Framework functions.

## Features

- **Guided Mapping**: Interactive prompts with definitions for Asset Classes and NIST Functions.
- **CSF 2.0 Coverage**: Includes the `Govern` function alongside Identify, Protect, Detect, Respond, and Recover.
- **CDM Principles**: Built-in support for "Left/Right of Boom" and "Applications vs. Devices" logic.
- **Dependency Mapping**: Capture Technology, People, and Process for every control.
- **Smart Resume**: Load existing JSON data and jump directly to specific assets or instances.
- **Grid Visualization**: Instant terminal-based coverage summary with completion percentages.
- **Quick-Copy Shortcut**: Duplicate mappings from similar assets to speed up data entry.
- **Multi-Format Export**: Generates styled Excel reports and raw CSV/JSON data.

## Requirements

- [Go](https://go.dev/) 1.26.1+
- [Make](https://www.gnu.org/software/make/) (optional, for using the Makefile)

## Usage

### Run the Wizard

The easiest way to start a new assessment is using `make`:

```bash
make run
```

Or using the Go CLI:

```bash
go run ./cmd/cdmbuddy
```

![sample-run](wizard.png)

### Command Line Options

CDM Buddy supports several flags for automation and resuming work:

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--input` | `-i` | Load CDM data from a JSON file to resume or report |
| `--output` | `-o` | Output Excel file path (default: `cdm_output.xlsx`) |
| `--csv` | `-c` | Output CSV file path (default: `cdm_output.csv`) |
| `--export-json` | | Export raw data as JSON (default: `cdm_data.json`) |
| `--name` | `-n` | Assessment name (used for default filenames) |
| `--report` | `-r` | Generate reports and exit without showing the wizard (requires `-i`) |

**Example: Resuming an assessment**
```bash
go run ./cmd/cdmbuddy -i my_assessment.json
```

### Run Tests

To verify the exporter and data structure:

```bash
make test
```

## Source Code

See [GUIDE.md](GUIDE.md) and [WORKFLOW.md](WORKFLOW.md) for details on the project structure and CDM alignment.

## Credits & Acknowledgments

This tool is an unofficial companion to the **Cyber Defense Matrix (CDM)** created by **Sounil Yu**. 

- **Official Website**: [cyberdefensematrix.com](https://cyberdefensematrix.com/)
- **The Book**: [*Cyber Defense Matrix: The Essential Guide to Navigating the Cybersecurity Landscape*](https://www.amazon.com/Cyber-Defense-Matrix-Navigating-Cybersecurity/dp/B09QP2GSGZ/) by Sounil Yu.

Inspired by [Stephen Dyson](https://www.linkedin.com/in/stephen-dyson-cybersecurity/)'s talk at [BSides Charm 2026](https://www.bsidescharm.org/schedule/).

The CDM is a trademark of Sounil Yu. This project is not affiliated with, endorsed by, or sponsored by Sounil Yu or the Cyber Defense Matrix organization.
