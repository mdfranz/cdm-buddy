# CDM Wizard Source

This directory contains the core logic for the Cyber Defense Matrix (CDM) Wizard.

## Structure

- `cdm_wizard/`: Main package.
    - `wizard.py`: Handles the interactive CLI flow, providing guided definitions for Asset Classes and Functions.
    - `exporter.py`: Logic for generating a styled Excel report that visualizes the "Left/Right of Boom" divide and the Technology-People-Process dependency continuum.

## CDM Alignment

The code follows Sounil Yu's *Cyber Defense Matrix* principles:
- **Asset Classes**: Devices (including commodity software), Networks, Applications (in-house built), Data, and Users.
- **Functions**: Govern, Identify, Protect, Detect, Respond, Recover.
- **Dichotomy**: Clearly separates "Left of Boom" from "Right of Boom".
- **Continuum**: Maps Technology, People, and Process for every cell, acknowledging the shifting dependency curve.
