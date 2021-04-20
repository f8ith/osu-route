import sys
import argparse

import requests
import webbrowser
import urllib.parse

from PySide6.QtWidgets import QWidget, QApplication, QMessageBox


def main():
    app = QApplication()
    parser = argparse.ArgumentParser("osu-route")
    parser.add_argument("link", type=str)
    args = parser.parse_args()
    if args.link.startswith("https://osu.ppy.sh/beatmaps/"):
        qm = QMessageBox()
        ret = qm.question(QWidget(parent=None), "", "Download beatmap?", qm.Yes | qm.No)

        if ret == qm.Yes:
            beatmap = args.link.split("/")[4]
            r = requests.get(f"https://chimu.moe/d/{beatmap}", stream=True)
            d = r.headers["Content-Disposition"]
            filename = (
                urllib.parse.unquote(d.split('filename="')[1].split('";')[0])
                .replace("/", "_")
                .replace('"', "")
                .replace("*", " ")
            )
            print(filename)
            with open(f"{filename}", "wb") as f:
                for chunk in r.iter_content(4096):
                    f.write(chunk)
            return
    webbrowser.open(args.link)


if __name__ == "__main__":
    main()
