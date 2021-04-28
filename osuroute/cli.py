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
    osu_path = Path("/mnt/HDD/osu")
    firefox_path = Path("/usr/bin/firefox")
    osu = "osu"
else:
    osu_path = Path("d:/osu/")
    firefox_path = "c:/Program Files/Mozilla Firefox/firefox.exe"
    osu = osu_path / "osu!.exe"


def main():
    parser = argparse.ArgumentParser("osu-route")
    parser.add_argument("link", type=str)
    args = parser.parse_args()
    processes = psutil.process_iter()
    ret = "osu!.exe" in (p.name() for p in processes)
    if ret:
        if args.link.startswith("https://osu.ppy.sh/b/") or args.link.startswith(
            "https://osu.ppy.sh/beatmaps/"
        ):
            beatmap_id = args.link.split("/")[4]
            beatmap = requests.get(f"https://api.chimu.moe/v1/map/{beatmap_id}")
            set_id = beatmap.json()["data"]["ParentSetId"]
        elif args.link.startswith("https://osu.ppy.sh/beatmapsets/"):
            set_id = args.link.split("/")[4]
        else:
            subprocess.run([Path(firefox_path), "-osint", "-url", args.link])
            return
        processes = psutil.process_iter()
        ret = "osu!.exe" in (p.name() for p in processes)
        if ret:
            r = requests.get(f"https://chimu.moe/d/{set_id}", stream=True)
            d = r.headers["Content-Disposition"]
            filename = (
                urllib.parse.unquote(d.split('filename="')[1].split('";')[0])
                .replace("/", "_")
                .replace('"', "")
                .replace("*", " ")
            )
            filepath = osu_path / filename
            with open(filepath, "wb") as f:
                for chunk in r.iter_content(4096):
                    f.write(chunk)
            subprocess.run([osu, filepath])
            return
    subprocess.run([Path(firefox_path), "-osint", "-url", args.link])


if __name__ == "__main__":
    main()
