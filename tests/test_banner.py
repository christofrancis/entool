"""Check that the README logo is printed once before the interactive menu."""

import contextlib
import io
import unittest
from unittest.mock import patch

from EnTool import BANNER, main


class StartupBannerTests(unittest.TestCase):
    def test_banner_is_printed_before_menu_and_only_once(self):
        output = io.StringIO()
        with patch("builtins.input", side_effect=["invalid", "0"]):
            with contextlib.redirect_stdout(output):
                main()

        text = output.getvalue()
        self.assertTrue(text.startswith(BANNER + "\n\n========= EnTool ========="))
        self.assertEqual(text.count("by christofrancis"), 1)
        self.assertIn("1. Encrypt a secret into a new text file", text)
        self.assertIn("2. Decrypt a secret from a text file", text)
        self.assertIn("Goodbye!", text)


if __name__ == "__main__":
    unittest.main()
