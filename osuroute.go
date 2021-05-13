package main

import (
	"runtime"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

type Config struct {
	OsuDir      string
	BrowserPath string
	Osu         string
}

var config Config

var SetId string

func LoadConfig() {
	var HomeDir string
	if runtime.GOOS != "windows" {
		HomeDir = path.Join(os.Getenv("HOME"), ".config/osuroute/")
	} else {
		HomeDir = os.Getenv("%LOCALAPPDATA%")
	}
	ConfigData, err := ioutil.ReadFile(path.Join(HomeDir, "osuroute.json"))
	if err != nil {
		log.Fatal("loadConfig:", err.Error())
	}
	err = json.Unmarshal(ConfigData, &config)
	if err != nil {
		log.Fatal("loadConfig:", err.Error())
	}
}

func OsuRunning() bool {
	processes, _ := process.Processes()
	for _, p := range processes {
		ProcessName, _ := p.Name()
		if ProcessName == "Osu!.exe" {
			return true
		}
	}
	return false
}

func main() {
	LoadConfig()
	OsuDir := config.OsuDir
	BrowserPath := config.BrowserPath
	Osu := config.Osu
	if len(os.Args) <= 1 {
		return
	}
	url := os.Args[1]
	if OsuRunning() {
		if strings.HasPrefix(url, "https://osu.ppy.sh/b/") || strings.HasPrefix(
			url, "https://osu.ppy.sh/beatmaps/") {
			BeatmapId := strings.Split(url, "/")[4]
			r, err := http.Get(fmt.Sprintf("https://api.chimu.moe/v1/map/%v", BeatmapId))
			if err != nil {
				log.Fatal(err)
			}
			defer r.Body.Close()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			var BeatmapData map[string]interface{}
			json.Unmarshal([]byte(body), &BeatmapData)
			data, ok := BeatmapData["data"].([]interface{})
			if !ok {
				return
			}
			SetId = data[1].(string)
		} else if strings.HasPrefix(url, "https://osu.ppy.sh/beatmapsets/") {
			SetId = strings.Split(url, "/")[4]
		} else {
			err := exec.Command(BrowserPath, url).Run()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		r, err := http.Get(fmt.Sprintf("https://api.chimu.moe/v1/download/%v?n=0", SetId))
		if err != nil {
			log.Fatal(err)
		}
		d := r.Header["Content-Disposition"][0]
		replacer := strings.NewReplacer("/", "_", `"`, "", "*", " ", "..", ".")
		filename := replacer.Replace(strings.Split(strings.Split(d, `filename="`)[1], `";`)[0])
		filepath := path.Join(OsuDir, filename)
		f, err := os.Create(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		io.Copy(f, r.Body)
		err = exec.Command(Osu, filepath).Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	err := exec.Command(BrowserPath, url).Run()
	if err != nil {
		log.Fatal(err)
	}
	return
}
