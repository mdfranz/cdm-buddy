import json
import questionary
from rich.console import Console
from rich.table import Table
from rich.panel import Panel

from cdm_wizard.model import (
    ASSET_DESCRIPTIONS,
    DATA_FIELDS,
    FUNCTION_DESCRIPTIONS,
    GOVERN_FUNCTIONS,
    RIGHT_OF_BOOM_FUNCTIONS,
    asset_classes,
    empty_matrix,
    functions,
)

console = Console()

def load_from_json(path: str) -> dict:
    """
    Load CDM data from a JSON file.

    Expected structure (unknown assets/functions are ignored; missing fields
    default to empty strings):

        {
          "Devices": {
            "Identify": {"Tech": "CrowdStrike", "People": "SysAdmin", "Process": "Patch SOP"},
            ...
          },
          ...
        }

    Only the canonical asset classes and NIST functions are accepted; any extra
    keys in the JSON are silently dropped so test fixtures can include comments
    or metadata under non-CDM keys without breaking the export.
    """
    assets = asset_classes()
    function_list = functions()

    with open(path) as f:
        raw = json.load(f)

    if not isinstance(raw, dict):
        raise ValueError(f"JSON root must be an object, got {type(raw).__name__}")

    # Build a clean data structure with defaults, accepting only known keys.
    data = empty_matrix()

    unknown_assets = [k for k in raw if k not in assets]
    if unknown_assets:
        console.print(f"[yellow]Warning:[/yellow] Ignoring unknown asset keys: {unknown_assets}")

    for asset in assets:
        if asset not in raw:
            continue
        asset_data = raw[asset]
        if not isinstance(asset_data, dict):
            raise ValueError(f"Asset '{asset}' must be an object, got {type(asset_data).__name__}")

        unknown_funcs = [k for k in asset_data if k not in function_list]
        if unknown_funcs:
            console.print(f"[yellow]Warning:[/yellow] {asset}: ignoring unknown function keys: {unknown_funcs}")

        for func in function_list:
            if func not in asset_data:
                continue
            cell = asset_data[func]
            if not isinstance(cell, dict):
                raise ValueError(f"Cell '{asset}.{func}' must be an object, got {type(cell).__name__}")
            data[asset][func] = {field: str(cell.get(field, "")) for field in DATA_FIELDS}

    return data


def run_wizard():
    """
    Runs the interactive CLI wizard to collect CDM data.
    """
    assets = asset_classes()
    function_list = functions()
    data = empty_matrix()

    console.print(Panel(
        "[bold blue]Cyber Defense Matrix Wizard[/bold blue]\n"
        "[italic]Guided by Sounil Yu's 'The Essential Guide to Navigating the Cybersecurity Landscape'[/italic]\n\n"
        "The CDM is a [bold]Mutually Exclusive and Collectively Exhaustive (MECE)[/bold] framework \n"
        "designed to provide structural clarity and identify gaps in your security posture."
    ))

    console.print(
        "\n[bold yellow]Core Concepts:[/bold yellow]\n"
        "• [bold blue]Left of Boom (Proactive):[/bold blue] [bold]Identify[/bold] and [bold]Protect[/bold]. Focuses on [italic]Structural Awareness[/italic]\n"
        "  (inventory, vulnerability discovery, and preventative safeguards).\n"
        "• [bold red]Right of Boom (Reactive):[/bold red] [bold]Detect[/bold], [bold]Respond[/bold], and [bold]Recover[/bold]. Focuses on [italic]Situational Awareness[/italic]\n"
        "  (identifying exploitations, containment, and restoration).\n"
        "\n[bold yellow]The Dependency Continuum:[/bold yellow]\n"
        "• [bold]Technology[/bold] dependency is strongest in [blue]Identify/Protect[/blue].\n"
        "• [bold]People[/bold] dependency grows significantly in [red]Detect/Respond/Recover[/red].\n"
        "• [bold]Process[/bold] is the essential, consistent foundation throughout all functions.\n"
    )

    from cdm_wizard.model import ASSET_ICONS

    for asset in assets:
        icon = ASSET_ICONS.get(asset, "🔹")
        console.print(f"\n[bold cyan]┏━ {icon} {asset.upper()} " + ("━" * (console.width - len(asset) - 10)) + "[/bold cyan]")
        console.print(f"  [italic dim cyan]{ASSET_DESCRIPTIONS[asset]}[/italic dim cyan]")

        for func in function_list:
            if func in GOVERN_FUNCTIONS:
                boom_tag = "[bold yellow]Cross-cutting[/bold yellow]"
                tech_instruction = "e.g. GRC platform, Policy Portal"
                people_instruction = "e.g. Security Governance Lead, Risk Committee"
                process_instruction = "e.g. Risk Exception Process, Security Policy Review"
            else:
                is_right_of_boom = func in RIGHT_OF_BOOM_FUNCTIONS
                boom_tag = "[bold red]Right of Boom[/bold red]" if is_right_of_boom else "[bold blue]Left of Boom[/bold blue]"
                tech_instruction = "e.g. EPP, Firewall, SIEM"
                people_instruction = "e.g. SOC Analyst, SysAdmin"
                process_instruction = "e.g. Patch Management SOP, IR Plan"

            console.print(f"\n  [bold reverse] {func} [/bold reverse] [bold]{boom_tag}[/bold]")
            console.print(f"  [dim]{FUNCTION_DESCRIPTIONS[func]}[/dim]")

            # Context-specific mapping tips
            tip_prefix = "  [italic yellow]💡 Tip:[/italic yellow]"
            if func == "Govern":
                console.print(f"{tip_prefix} [italic dim]Think of this as the 'Context'—policies, risk appetite, and oversight governing this asset class.[/italic dim]")
            elif func in ("Identify", "Protect", "Recover"):
                console.print(f"{tip_prefix} [italic dim]Map based on the primary asset being acted upon.[/italic dim]")
            elif func == "Detect":
                console.print(f"{tip_prefix} [italic dim]Map based on the Use Case (e.g., Insider Threat maps to Users).[/italic dim]")
            elif func == "Respond":
                console.print(f"{tip_prefix} [italic dim]Map based on the asset being investigated or contained.[/italic dim]")

            while True:
                try:
                    tech = questionary.text(
                        f"    Technology/Vendor:",
                        instruction=tech_instruction,
                        qmark=""
                    ).unsafe_ask()
                    data[asset][func]["Tech"] = tech if tech else ""

                    people = questionary.text(
                        f"    People/Responsible Role:",
                        instruction=people_instruction,
                        qmark=""
                    ).unsafe_ask()
                    data[asset][func]["People"] = people if people else ""

                    process = questionary.text(
                        f"    Process/Procedure:",
                        instruction=process_instruction,
                        qmark=""
                    ).unsafe_ask()
                    data[asset][func]["Process"] = process if process else ""
                    break

                except (KeyboardInterrupt, EOFError):
                    console.print("") # New line after the interrupt symbol
                    if questionary.confirm("Exit wizard and save progress?", qmark="", default=False).ask():
                        return data
                    console.print("[yellow]Resuming current function...[/yellow]")

    return data

def display_summary(data):
    """
    Displays a summary table of the collected data in the terminal.
    """
    function_list = functions()
    table = Table(title="Cyber Defense Matrix Summary", show_lines=True, expand=True)

    table.add_column("Asset", justify="left", style="bold cyan", no_wrap=True, min_width=12)
    for func in function_list:
        table.add_column(func, style="magenta")

    for asset, mappings in data.items():
        row_cells = [asset]
        for func in function_list:
            cell_data = mappings.get(func, {"Tech": "", "People": "", "Process": ""})
            tech    = cell_data["Tech"]    or "[dim]—[/dim]"
            people  = cell_data["People"]  or "[dim]—[/dim]"
            process = cell_data["Process"] or "[dim]—[/dim]"
            summary = f"[bold]Tech:[/bold] {tech}\n[bold]People:[/bold] {people}\n[bold]Process:[/bold] {process}"
            row_cells.append(summary)
        table.add_row(*row_cells)

    console.print(table)
