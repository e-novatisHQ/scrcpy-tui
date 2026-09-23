#!/usr/bin/env python3
"""Isolated UX-01 proposal; no imports or mutations of the shipped launcher."""
import argparse
import os
import string

CATALOGS = {
    "fr": {
        "ready": "Prêt à lancer · Entrée",
        "devices": "Appareils",
        "presets": "Presets",
        "unauthorized": "Non autorisé : {device}",
        "keys": "Entrée lancer · r rafraîchir · q quitter",
    },
    "en": {
        "ready": "Ready to launch · Enter",
        "devices": "Devices",
        "presets": "Presets",
        "unauthorized": "Unauthorized: {device}",
        "keys": "Enter launch · r refresh · q quit",
    },
}


def normalize(value):
    return value.strip().split(".", 1)[0].split("@", 1)[0].replace("_", "-").split("-", 1)[0].lower()


def select_locale(explicit, environment):
    if explicit is not None:
        candidate = normalize(explicit)
        if candidate not in CATALOGS:
            raise ValueError("Locale attendue : fr ou en")
        return candidate
    for key in ("LC_ALL", "LC_MESSAGES", "LANG"):
        if environment.get(key):
            candidate = normalize(environment[key])
            return candidate if candidate in CATALOGS else "fr"
    return "fr"


def self_test():
    cases = [
        (None, {}, "fr"),
        ("en-US", {"LC_ALL": "fr"}, "en"),
        ("fr_FR.UTF-8", {"LANG": "en"}, "fr"),
        (None, {"LC_ALL": "en_GB.UTF-8", "LC_MESSAGES": "fr"}, "en"),
        (None, {"LC_MESSAGES": "en", "LANG": "fr"}, "en"),
        (None, {"LANG": "en_US.UTF-8"}, "en"),
        (None, {"LC_ALL": "de_DE", "LANG": "en"}, "fr"),
        (None, {"LANG": "C.UTF-8"}, "fr"),
        (None, {"LANG": "POSIX"}, "fr"),
    ]
    for explicit, environment, expected in cases:
        assert select_locale(explicit, environment) == expected
    for invalid in ("", "de", "zz"):
        try:
            select_locale(invalid, {})
        except ValueError:
            pass
        else:
            raise AssertionError("Unknown explicit locale accepted")
    def fields(text):
        return [field for _, field, _, _ in string.Formatter().parse(text) if field is not None]
    assert CATALOGS["fr"].keys() == CATALOGS["en"].keys()
    for key in CATALOGS["fr"]:
        assert fields(CATALOGS["fr"][key]) == fields(CATALOGS["en"][key])
        assert all(CATALOGS[locale][key] for locale in CATALOGS)
    print("Locale precedence, fallback, catalog parity and placeholders: PASS")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--locale")
    parser.add_argument("--self-test", action="store_true")
    args = parser.parse_args()
    if args.self_test:
        self_test()
        return
    try:
        locale = select_locale(args.locale, os.environ)
    except ValueError as error:
        parser.error(str(error))
    text = CATALOGS[locale]
    print(f"scrcpy-tui — prototype {locale.upper()}")
    print(text["ready"])
    print(f'{text["devices"]:<34}{text["presets"]}')
    print("> Téléviseur de démonstration      > Léger Wi-Fi")
    print("  " + text["unauthorized"].format(device="téléphone test"))
    print(text["keys"])


if __name__ == "__main__":
    main()
