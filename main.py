import sys
import os
import argparse

# Add src to path if needed for local runs
sys.path.append(os.path.join(os.path.dirname(__file__), 'src'))

from cdm_wizard.wizard import run_wizard, display_summary, load_from_json
from cdm_wizard.exporter import export_to_excel
from rich.console import Console

console = Console()
error_console = Console(stderr=True)

def main(argv=None):
    parser = argparse.ArgumentParser(description="Cyber Defense Matrix Wizard")
    parser.add_argument(
        "--input", "-i",
        metavar="FILE",
        help="Load CDM data from a JSON file instead of running the interactive wizard",
    )
    parser.add_argument(
        "--output", "-o",
        metavar="FILE",
        default="cdm_output.xlsx",
        help="Output Excel file path (default: cdm_output.xlsx)",
    )
    args = parser.parse_args(argv)

    try:
        if args.input:
            console.print(f"[bold cyan]Loading data from:[/bold cyan] {args.input}")
            data = load_from_json(args.input)
        else:
            data = run_wizard()

        display_summary(data)

        export_to_excel(data, args.output)
        console.print(f"\n[bold green]Success![/bold green] Your Cyber Defense Matrix has been exported to [bold cyan]{args.output}[/bold cyan]")
        return 0

    except KeyboardInterrupt:
        error_console.print("\n[yellow]Wizard cancelled by user.[/yellow]")
        return 130
    except Exception as e:
        error_console.print(f"\n[bold red]Error:[/bold red] {str(e)}")
        return 1

if __name__ == "__main__":
    raise SystemExit(main())
