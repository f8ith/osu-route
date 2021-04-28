import sys
import platform
import argparse
import subprocess
from pathlib import Path

import requests
import webbrowser
import urllib.parse

import psutil

if platform.system() == "Linux":
    osu_path = "/mnt/HDD/osu"
    firefox_path = "/usr/bin/firefox"
else:
    osu_path = "c:/Users/jamie/"
    firefox_path = "c:/Program Files/Mozilla Firefox/firefox.exe"


def main():
    parser = argparse.ArgumentParser("osu-route")
    parser.add_argument("link", type=str)
    args = parser.parse_args()
    if args.link.startswith("https://osu.ppy.sh/b/"):
        processes = psutil.process_iter()
        ret = "osu!.exe" in (p.name() for p in processes)
        if ret:
            beatmap = args.link.split("/")[4]
            r = requests.get(f"https://chimu.moe/d/{beatmap}", stream=True)
            d = r.headers["Content-Disposition"]
            filename = (
                urllib.parse.unquote(d.split('filename="')[1].split('";')[0])
                .replace("/", "_")
                .replace('"', "")
                .replace("*", " ")
            )
            with open(Path(osu_path) / filename, "wb") as f:
                for chunk in r.iter_content(4096):
                    f.write(chunk)
            return
    subprocess.run([Path(firefox_path), "-osint", "-url", args.link])


if __name__ == "__main__":
    main()
