# Cyber Defense Matrix: Initial Population Workflow

This document outlines the ideal workflow for the initial population and operationalization of the Cyber Defense Matrix (CDM).

## **Core Mapping Principles**
Before beginning, adhere to these "Rules of Thumb" to ensure consistency:
*   **The First-Order Rule:** Map tools based on their primary function, not where they live or their secondary benefits (e.g., a WAF lives on the Network but primarily **Protects Applications**).
*   **The Verb Rule:** Focus on what the tool *does* (e.g., "Does it prevent?" = Protect; "Does it notify?" = Detect).
*   **Functional Consistency:** Ensure the same sub-function (like "Inventory") is applied consistently across all five asset classes (Devices, Networks, Apps, Data, Users).
*   **Dependency Continuum:** Remember that "Left of Boom" (Identify/Protect) is technology-heavy, while "Right of Boom" (Detect, Respond, Recover) is people-heavy.

---

## **Phase 1: Inventory Discovery (The "Identify" Foundation)**
The goal is to achieve structural awareness of the environment.
1.  **Enumerate Assets:** Gather existing inventories for all five asset classes.
2.  **Identify Criticality:** Note "Crown Jewel" assets to prioritize mapping.
3.  **Review Governance:** Identify policies, compliance frameworks (NIST CSF, CIS), and existing risk assessments.

## **Phase 2: Mapping Exercise (The "Action" Phase)**
Systematically populate the 5x5 grid with your current security tech stack and processes.
1.  **Map Primary Functions:** Place each tool in its most relevant cell. 
2.  **Deconstruct Suites:** For platforms that span multiple cells (e.g., EDR/XDR), split their specific features into the respective boxes (e.g., EDR Alerting -> Detect/Device; EDR Quarantine -> Respond/Device).
3.  **Account for Non-Tech:** Map human-centric processes (e.g., "Manual Log Review") into the "Right of Boom" cells where technology dependency is lower.

## **Phase 3: Gap & Overlap Analysis (The "Insight" Phase)**
Analyze the grid to identify strategic imbalances.
1.  **Identify Blind Spots:** Mark empty cells. These are areas where the organization lacks defense-in-depth or visibility.
2.  **Identify Redundancies:** Flag cells with multiple overlapping tools. These are prime candidates for vendor consolidation and cost savings.
3.  **Evaluate Dependency:** Ensure that "Right of Boom" functions have adequate human resources (People) and documented procedures (Process) to act on technology alerts.

## **Phase 4: Strategy & Prioritization**
1.  **Fill Critical Gaps:** Prioritize investment in empty cells that represent high-risk attack surfaces.
2.  **Rationalize Spend:** Use overlap data to justify retiring redundant tools and reallocating budget.
3.  **Align with Maturity:** Focus on stabilizing "Left of Boom" (Identify/Protect) before attempting sophisticated "Right of Boom" (Detect/Respond) automation.

## **Phase 5: Operationalization**
1.  **Executive Communication:** Use a visual template (Slides) to present the matrix to leadership, framing technical gaps as business risks.
2.  **Continuous Lifecycle:** Update the matrix during:
    *   New tool procurement.
    *   Major architectural changes (e.g., Cloud migration).
    *   Annual risk assessments.
3.  **Standardize Terminology:** Use the CDM to normalize language across IT, Security, and Audit teams.
