import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[1]
FIXTURE_PATH = REPO_ROOT / "tests" / "fixtures" / "sample_cdm.json"


class CLITests(unittest.TestCase):
    def run_cli(self, *args: str) -> subprocess.CompletedProcess[str]:
        return subprocess.run(
            [sys.executable, "main.py", *args],
            cwd=REPO_ROOT,
            capture_output=True,
            text=True,
        )

    def test_missing_input_returns_nonzero_exit_code(self):
        result = self.run_cli("--input", "does-not-exist.json")

        self.assertNotEqual(result.returncode, 0)
        self.assertIn("Error:", result.stderr)

    def test_invalid_json_root_returns_nonzero_exit_code(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            invalid_json_path = Path(tmpdir) / "invalid.json"
            invalid_json_path.write_text("[]", encoding="utf-8")

            result = self.run_cli("--input", str(invalid_json_path))

        self.assertNotEqual(result.returncode, 0)
        self.assertIn("JSON root must be an object", result.stderr)

    def test_successful_json_import_returns_zero_and_writes_output(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            output_path = Path(tmpdir) / "cdm_output.xlsx"

            result = self.run_cli(
                "--input", str(FIXTURE_PATH),
                "--output", str(output_path),
            )

            self.assertEqual(result.returncode, 0, msg=result.stderr)
            self.assertTrue(output_path.exists())
            self.assertIn("Success!", result.stdout)


if __name__ == "__main__":
    unittest.main()
