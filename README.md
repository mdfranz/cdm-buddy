![CDM Buddy Logo](cdm-buddy.png)

# Cyber Defense Matrix (CDM) Helper Monkey

An interactive CLI tool to map your security portfolio to Sounil Yu's Cyber Defense Matrix.


`NOTE`: This is a work in progress as I'm still trying to figure out how I'll use CDM.

## Features

- **Guided Mapping**: Interactive prompts with definitions for Asset Classes and NIST Functions.
- **CSF 2.0 Coverage**: Includes the `Govern` function alongside Identify, Protect, Detect, Respond, and Recover.
- **Six Asset Classes**: Coverage for Devices, Networks, Applications, Data, Users, and the **newly added Services (SaaS)** class.
- **CDM Principles**: Built-in support for "Left/Right of Boom" and "Applications vs. Devices" logic.
- **Dependency Mapping**: Capture Technology, People, and Process for every control.
- **Multi-Format Export**: Generates styled **Excel** reports, **CSV** raw data, and **JSON** for session persistence.
- **Resume & Edit**: Load existing JSON data to continue an assessment or edit previous entries.

## Requirements

- [Go](https://go.dev/) 1.26.1+
- [Make](https://www.gnu.org/software/make/) (optional, for using the Makefile)

## Usage

### Run the Wizard

The easiest way to run the wizard is using `make`:

```bash
make run
```

Or using the Go CLI:

```bash
cdmbuddy
```

Example with sample data

```
mdfranz@lenovo-cr14p-arm:~/github/cdm-buddy$ ./cdmbuddy -i samples/sample.json -r
Loading data from: samples/sample.json
                                 
 📊 Cyber Defense Matrix Coverage
                                 
  ┌───────────────────────────────┬────────────┬────────────┬────────────┬────────────┬────────────┬────────────┐
  │ Instance                      │   Govern   │  Identify  │  Protect   │   Detect   │  Respond   │  Recover   │
  ├───────────────────────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┤
  │ Microsoft 365                 │     ✓      │     ─      │     ✓      │     ─      │     ─      │     ─      │
  ├───────────────────────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┤
  │ Corporate Laptops             │     ─      │     ✓      │     ✓      │     ✓      │     ✓      │     ✓      │
  ├───────────────────────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┤
  │ Legacy Servers                │     ─      │     ✓      │     ❌     │     ✓      │     ─      │     ─      │
  ├───────────────────────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┤
  │ Office WIFI                   │     ─      │     ─      │     ✓      │     ✓      │     ─      │     ─      │
  ├───────────────────────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┤
  │ Customer Portal               │     ✓      │     ✓      │     ✓      │     ─      │     ─      │     ─      │
  ├───────────────────────────────┼────────────┼────────────┼────────────┼────────────┼────────────┼────────────┤
  │ All Employees                 │     ─      │     ✓      │     ✓      │     ─      │     ─      │     ─      │
  └───────────────────────────────┴────────────┴────────────┴────────────┴────────────┴────────────┴────────────┘
                                                                         
╭───────────────────────────────────────────────────────────────────────╮
│    📈   16/36 cells filled  •  44% coverage  •  1 marked no coverage  │
╰───────────────────────────────────────────────────────────────────────╯

Success! Your Cyber Defense Matrix has been exported to cdm_output.xlsx

```

### CLI Flags

| Flag | Shorthand | Description | Default |
|------|-----------|-------------|---------|
| `--input` | `-i` | Load CDM data from a JSON file to resume | |
| `--output` | `-o` | Output Excel file path | `cdm_output.xlsx` |
| `--csv` | `-c` | Output CSV file path | `cdm_output.csv` |
| `--export-json`| | Also export CDM data as JSON | `cdm_data.json` |
| `--name` | `-n` | Assessment name (used for filenames) | |
| `--report` | `-r` | Generate report and exit (requires `-i`) | `false` |

### Example: Resume an Assessment
```bash
cdmbuddy -i my_previous_run.json
```

### Run Tests

To verify the exporter and data structure:

```bash
make test
```

Or:

```bash
go test -v ./...
```

## Source Code

See [GUIDE.md](GUIDE.md) and [WORKFLOW.md](WORKFLOW.md) for details on the project structure and CDM alignment.

## Credits & Acknowledgments

This tool is an unofficial companion to the **Cyber Defense Matrix (CDM)** created by **Sounil Yu**. 

- **Official Website**: [cyberdefensematrix.com](https://cyberdefensematrix.com/)
- **The Book**: [*Cyber Defense Matrix: The Essential Guide to Navigating the Cybersecurity Landscape*](https://www.amazon.com/Cyber-Defense-Matrix-Navigating-Cybersecurity/dp/B09QP2GSGZ/) by Sounil Yu.

I got inspired to write this after seeing a talk by [Stephen Dyson](https://www.linkedin.com/in/stephen-dyson-cybersecurity/) at [Bsides Charm 2026](https://www.bsidescharm.org/schedule/)

The CDM is a trademark of Sounil Yu. This project is not affiliated with, endorsed by, or sponsored by Sounil Yu or the Cyber Defense Matrix organization.
