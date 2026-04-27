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

TECH_EXAMPLES = {
    ("Devices", "Govern"):      "e.g. CMDB policies, device standards documentation",
    ("Devices", "Identify"):    "e.g. Qualys, Tenable, Rapid7, Nessus",
    ("Devices", "Protect"):     "e.g. CrowdStrike Falcon, SentinelOne, Tanium, Microsoft Defender",
    ("Devices", "Detect"):      "e.g. EDR telemetry, Sysmon, Microsoft Defender for Endpoint",
    ("Devices", "Respond"):     "e.g. CrowdStrike RTR, Tanium, Carbon Black Live Response",
    ("Devices", "Recover"):     "e.g. Veeam, Cohesity, OS reimaging pipeline",
    ("Networks", "Govern"):     "e.g. Network segmentation policy, firewall rule change management",
    ("Networks", "Identify"):   "e.g. Nmap, Shodan, Censys, network asset scanners",
    ("Networks", "Protect"):    "e.g. Palo Alto NGFW, Cisco ASA, Zscaler, Cloudflare",
    ("Networks", "Detect"):     "e.g. Zeek/Bro, Suricata, Darktrace, Cisco Stealthwatch",
    ("Networks", "Respond"):    "e.g. Firewall ACL scripts, network isolation, Cisco ISE",
    ("Networks", "Recover"):    "e.g. SD-WAN failover, BGP route restoration, ISP failover",
    ("Applications", "Govern"):  "e.g. SDLC policies, AppSec program charter, risk acceptance process",
    ("Applications", "Identify"): "e.g. Snyk, OWASP ZAP, Veracode, Semgrep, SBOM tools",
    ("Applications", "Protect"): "e.g. WAF (Imperva, AWS WAF), Checkmarx, GitHub Advanced Security",
    ("Applications", "Detect"):  "e.g. RASP tools, AWS GuardDuty, application log monitoring",
    ("Applications", "Respond"): "e.g. Feature flags, hotfix pipeline, incident runbooks",
    ("Applications", "Recover"): "e.g. Blue/green deployment, rollback pipeline, DR runbooks",
    ("Data", "Govern"):         "e.g. Data classification policy, DLP policy, retention standards",
    ("Data", "Identify"):       "e.g. AWS Macie, Varonis, BigID, data catalog tools",
    ("Data", "Protect"):        "e.g. Varonis, Macie, Forcepoint DLP, encryption at rest",
    ("Data", "Detect"):         "e.g. Varonis DatAlert, CASB (Netskope, MCAS), database activity monitoring",
    ("Data", "Respond"):        "e.g. Data quarantine procedures, access revocation, legal hold",
    ("Data", "Recover"):        "e.g. S3 versioning, database backups, Cohesity DataProtect",
    ("Users", "Govern"):        "e.g. IAM policies, acceptable use policy, MFA requirements",
    ("Users", "Identify"):      "e.g. SailPoint, Saviynt, AD/LDAP inventory, Okta Lifecycle Management",
    ("Users", "Protect"):       "e.g. Okta, CyberArk, Duo, Microsoft Entra ID, BeyondTrust",
    ("Users", "Detect"):        "e.g. Okta ThreatInsight, Microsoft Sentinel UEBA, Exabeam",
    ("Users", "Respond"):       "e.g. Account disable scripts, CyberArk session termination, SOAR playbooks",
    ("Users", "Recover"):       "e.g. Account restoration procedures, password reset workflows",
}
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
