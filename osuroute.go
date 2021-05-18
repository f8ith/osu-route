package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/process"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	OsuDir      string
	BrowserPath string
	Osu         string
}

var config Config

var SetId string

var HomeDir string

var ConfigPath string

func LoadConfig() {
	HomeDir, _ = os.UserHomeDir()
	ConfigPath = path.Join(HomeDir, "osuroute.json")

	ConfigData, err := ioutil.ReadFile(ConfigPath)
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
		if ProcessName == "osu!.exe" {
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
		a := app.New()
		w := a.NewWindow("osuRoute!")

		entry1 := widget.NewEntry()
		entry1.Text = config.OsuDir
		entry2 := widget.NewEntry()
		entry2.Text = config.BrowserPath
		entry3 := widget.NewEntry()
		entry3.Text = config.Osu

		form := &widget.Form{
			Items: []*widget.FormItem{ // we can specify items in the constructor
				{Text: "Osu Directory", Widget: entry1},
				{Text: "Browser Path", Widget: entry2},
				{Text: "Osu App Path", Widget: entry3}},
			OnSubmit: func() { // optional, handle form submission
				config.OsuDir = entry1.Text
				config.BrowserPath = entry2.Text
				config.Osu = entry3.Text
				newConfig, _ := json.MarshalIndent(config, "", " ")
				ioutil.WriteFile(ConfigPath, newConfig, 0644)
				w.Close()
			},
			SubmitText: "Save",
		}
		w.SetContent(container.NewVBox(form))
		w.Resize(fyne.NewSize(450, 100))
		w.ShowAndRun()
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
