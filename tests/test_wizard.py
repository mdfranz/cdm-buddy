import json
import os
import sys
import tempfile
import unittest


sys.path.append(os.path.join(os.path.dirname(__file__), "../src"))

from cdm_wizard.model import empty_matrix
from cdm_wizard.wizard import console, display_summary, load_from_json


class WizardTests(unittest.TestCase):
    def write_json(self, payload) -> str:
        tmp = tempfile.NamedTemporaryFile("w", delete=False, encoding="utf-8")
        try:
            json.dump(payload, tmp)
            return tmp.name
        finally:
            tmp.close()

    def test_load_from_json_supports_govern_and_defaults_missing_fields(self):
        json_path = self.write_json(
            {
                "Devices": {
                    "Govern": {
                        "Tech": "ServiceNow GRC",
                        "People": "Security Governance",
                    },
                    "Identify": {
                        "Tech": "Lansweeper",
                    },
                },
                "UnknownAsset": {},
            }
        )
        try:
            with console.capture():
                data = load_from_json(json_path)
        finally:
            os.remove(json_path)

        self.assertEqual(data["Devices"]["Govern"]["Tech"], "ServiceNow GRC")
        self.assertEqual(data["Devices"]["Govern"]["People"], "Security Governance")
        self.assertEqual(data["Devices"]["Govern"]["Process"], "")
        self.assertEqual(data["Devices"]["Identify"]["Tech"], "Lansweeper")
        self.assertEqual(data["Devices"]["Identify"]["People"], "")
        self.assertNotIn("UnknownAsset", data)

    def test_load_from_json_rejects_non_object_root(self):
        json_path = self.write_json([])
        try:
            with self.assertRaisesRegex(ValueError, "JSON root must be an object"):
                load_from_json(json_path)
        finally:
            os.remove(json_path)

    def test_display_summary_includes_govern_column(self):
        data = empty_matrix()
        data["Devices"]["Govern"] = {
            "Tech": "Policy Portal",
            "People": "Security Governance",
            "Process": "Risk exception workflow",
        }

        with console.capture() as capture:
            display_summary(data)
        output = capture.get()

        self.assertIn("Govern", output)
        self.assertIn("Policy", output)
        self.assertIn("Portal", output)
        self.assertIn("Devices", output)


if __name__ == "__main__":
    unittest.main()
