ASSET_ICONS = {
    "Devices": "💻",
    "Networks": "🌐",
    "Applications": "📦",
    "Data": "💾",
    "Users": "👥",
}

ASSET_DESCRIPTIONS = {
    "Devices": "Workstations, servers, OS, firmware, and commodity software (e.g., MS Office, Browsers).",
    "Networks": "Communication channels, protocols (DNS, BGP), and paths (VPCs, VPNs). Not the physical hardware.",
    "Applications": "ONLY in-house built applications or software where you build/maintain the source code.",
    "Data": "Information at rest, in motion, or in use (Databases, S3 buckets, files).",
    "Users": "The people using the resources and their associated identities.",
}

FUNCTION_DESCRIPTIONS = {
    "Govern": "The 'Operating System' of the matrix. Defines strategy, policy, risk appetite, and oversight that guide all other functions.",
    "Identify": "Inventory, classification, and vulnerability discovery.",
    "Protect": "Safeguards to ensure delivery of services; patching and access controls.",
    "Detect": "Identifying the occurrence of a cybersecurity event or exploitation.",
    "Respond": "Activities to contain the impact of a detected event.",
    "Recover": "Resilience and restoration of capabilities impaired by an event.",
}

FUNCTION_COLORS = {
    "Govern": "D9EAD3",
    "Identify": "DCE6F1",
    "Protect": "E4DFEC",
    "Detect": "FFF2CC",
    "Respond": "FDE9D9",
    "Recover": "EBF1DE",
}

DATA_FIELDS = ("Tech", "People", "Process")
GOVERN_FUNCTIONS = ("Govern",)
LEFT_OF_BOOM_FUNCTIONS = ("Identify", "Protect")
RIGHT_OF_BOOM_FUNCTIONS = ("Detect", "Respond", "Recover")


def asset_classes() -> list[str]:
    return list(ASSET_DESCRIPTIONS.keys())


def functions() -> list[str]:
    return list(FUNCTION_DESCRIPTIONS.keys())


def blank_cell() -> dict[str, str]:
    return {field: "" for field in DATA_FIELDS}


def empty_matrix() -> dict[str, dict[str, dict[str, str]]]:
    return {
        asset: {func: blank_cell() for func in functions()}
        for asset in asset_classes()
    }
