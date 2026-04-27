package model

import (
	"encoding/json"
	"os"
)

var AssetIcons = map[string]string{
	"Devices":      "💻",
	"Networks":     "🌐",
	"Applications": "📦",
	"Data":         "💾",
	"Users":        "👥",
	"Services":     "☁️",
}

var AssetDescriptions = map[string]string{
	"Devices":      "Workstations, servers, OS, firmware, and commodity software (e.g., MS Office, Browsers, Linux).",
	"Networks":     "Communication channels, protocols (DNS, BGP), and paths (VPCs, VPNs). Not the physical hardware.",
	"Applications": "ONLY in-house built applications or software where you build/maintain the source code.",
	"Data":         "Information at rest, in motion, or in use (Databases, S3 buckets, files).",
	"Users":        "The people using the resources and their associated identities.",
	"Services":     "Third-party SaaS applications and cloud-delivered services (e.g., Salesforce, ServiceNow, Atlassian).",
}

var FunctionDescriptions = map[string]string{
	"Govern":   "The 'Operating System' of the matrix. Binds all functions through strategy, policy, and oversight to ensure alignment with organizational goals.",
	"Identify": "Foundational inventory and vulnerability discovery. High Technology dependency—you must know what you have before you can protect it.",
	"Protect":  "Proactive safeguards like patching and access controls. Technology-heavy focus on hardening assets and limiting the impact of threats.",
	"Detect":   "Real-time identification of security events. The transition from proactive to reactive, requiring a balance of Technology and People.",
	"Respond":  "Containment and mitigation once an event is detected. High People dependency; focused on stopping the 'bleeding' and coordination.",
	"Recover":  "Resilience and restoration of services. People-centric focus on returning to normal operations and learning from the incident.",
}

var FunctionColors = map[string]string{
	"Govern":   "D9EAD3",
	"Identify": "DCE6F1",
	"Protect":  "E4DFEC",
	"Detect":   "FFF2CC",
	"Respond":  "FDE9D9",
	"Recover":  "EBF1DE",
}

var TechExamples = map[string]string{
	"Devices-Govern":      "e.g. CMDB policies, device standards documentation",
	"Devices-Identify":    "e.g. Qualys, Tenable, Rapid7, Nessus, or Cloud Provider APIs",
	"Devices-Protect":     "e.g. CrowdStrike Falcon, SentinelOne, Tanium, Microsoft Defender",
	"Devices-Detect":      "e.g. EDR telemetry, Sysmon, Osquery, Auditd, Microsoft Defender for Endpoint",
	"Devices-Respond":     "e.g. CrowdStrike RTR, Tanium, Carbon Black Live Response",
	"Devices-Recover":     "e.g. Veeam, Cohesity, OS reimaging pipeline, or Cloud Provider automation",
	"Networks-Govern":     "e.g. Network segmentation policy, firewall rule change management",
	"Networks-Identify":   "e.g. Nmap, Shodan, Censys, network asset scanners",
	"Networks-Protect":    "e.g. Palo Alto NGFW, Cisco ASA, Zscaler, Cloudflare,",
	"Networks-Detect":     "e.g. Zeek/Bro, Suricata, Darktrace, Cisco Stealthwatch, Flow logging",
	"Networks-Respond":    "e.g. Firewall ACL scripts, network isolation, Cisco ISE",
	"Networks-Recover":    "e.g. SD-WAN failover, BGP route restoration, ISP failover",
	"Applications-Govern":  "e.g. SDLC policies, AppSec program charter, risk acceptance process",
	"Applications-Identify": "e.g. Snyk, OWASP ZAP, Veracode, Semgrep, SBOM tools",
	"Applications-Protect": "e.g. WAF (Imperva, AWS WAF), Checkmarx, GitHub Advanced Security",
	"Applications-Detect":  "e.g. RASP tools, AWS GuardDuty, application log monitoring",
	"Applications-Respond": "e.g. Feature flags, hotfix pipeline, incident runbooks",
	"Applications-Recover": "e.g. Blue/green deployment, rollback pipeline, DR runbooks",
	"Data-Govern":         "e.g. Data classification policy, DLP policy, retention standards",
	"Data-Identify":       "e.g. AWS Macie, Varonis, BigID, data catalog tools",
	"Data-Protect":        "e.g. Varonis, Macie, Forcepoint DLP, encryption at rest",
	"Data-Detect":         "e.g. Varonis DatAlert, CASB (Netskope, MCAS), database activity monitoring",
	"Data-Respond":        "e.g. Data quarantine procedures, access revocation, legal hold",
	"Data-Recover":        "e.g. S3 versioning, database backups, Cohesity DataProtect",
	"Users-Govern":        "e.g. IAM policies, acceptable use policy, MFA requirements",
	"Users-Identify":      "e.g. SailPoint, Saviynt, AD/LDAP inventory, Okta Lifecycle Management",
	"Users-Protect":       "e.g. Okta, CyberArk, Duo, Microsoft Entra ID, BeyondTrust",
	"Users-Detect":        "e.g. Okta ThreatInsight, Microsoft Sentinel UEBA, Exabeam",
	"Users-Respond":       "e.g. Account disable scripts, CyberArk session termination, SOAR playbooks",
	"Users-Recover":       "e.g. Account restoration procedures, password reset workflows",
}

var AssetClasses = []string{"Services", "Devices", "Networks", "Applications", "Data", "Users"}
var Functions = []string{"Govern", "Identify", "Protect", "Detect", "Respond", "Recover"}
var GovernFunctions = []string{"Govern"}
var LeftOfBoomFunctions = []string{"Identify", "Protect"}
var RightOfBoomFunctions = []string{"Detect", "Respond", "Recover"}

type Cell struct {
	Tech    string `json:"Tech"`
	People  string `json:"People"`
	Process string `json:"Process"`
}

type AssetInstance struct {
	Name  string          `json:"Name"`
	Cells map[string]Cell `json:"Cells"`
}

type Matrix map[string][]AssetInstance

func EmptyMatrix() Matrix {
	m := make(Matrix)
	for _, asset := range AssetClasses {
		m[asset] = []AssetInstance{}
	}
	return m
}

func LoadFromJson(path string) (Matrix, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	matrix := EmptyMatrix()
	if err := json.Unmarshal(data, &matrix); err != nil {
		return nil, err
	}

	return matrix, nil
}

func SaveToJson(matrix Matrix, path string) error {
	data, err := json.MarshalIndent(matrix, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func HasInstance(matrix Matrix, assetClass, name string) bool {
	instances := matrix[assetClass]
	for _, inst := range instances {
		if inst.Name == name {
			return true
		}
	}
	return false
}
