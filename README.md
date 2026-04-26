# Cyber Defense Matrix (CDM) Wizard

An interactive CLI tool to map your security portfolio to Sounil Yu's Cyber Defense Matrix.

## Features

- **Guided Mapping**: Interactive prompts with definitions for Asset Classes and NIST Functions.
- **CSF 2.0 Coverage**: Includes the `Govern` function alongside Identify, Protect, Detect, Respond, and Recover.
- **CDM Principles**: Built-in support for "Left/Right of Boom" and "Applications vs. Devices" logic.
- **Dependency Mapping**: Capture Technology, People, and Process for every control.
- **Excel Export**: Generates a styled report with visual cues for strategic analysis.

## Requirements

- [uv](https://github.com/astral-sh/uv) (recommended)
- Python 3.14+

## Usage

### Run the Wizard

The easiest way to run the wizard is using `uv`:

```bash
cd cdm-wizard
uv run main.py
```

### Run Tests

To verify the exporter and data structure:

```bash
cd cdm-wizard
uv run -m unittest discover -s tests
```

## Source Code

See [src/README.md](src/README.md) for details on the internal package structure and CDM alignment.

## Credits & Acknowledgments

This tool is an unofficial companion to the **Cyber Defense Matrix (CDM)** created by **Sounil Yu**. 

- **Official Website**: [cyberdefensematrix.com](https://cyberdefensematrix.com/)
- **The Book**: [*Cyber Defense Matrix: The Essential Guide to Navigating the Cybersecurity Landscape*](https://www.amazon.com/Cyber-Defense-Matrix-Essential-Cybersecurity/dp/1955419311) by Sounil Yu.

The CDM is a trademark of Sounil Yu. This project is not affiliated with, endorsed by, or sponsored by Sounil Yu or the Cyber Defense Matrix organization.
