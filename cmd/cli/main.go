package main

import (
	"log"
	"math"
	"os"
	"regexp"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type nodeValue string

func (nv nodeValue) String() string {
	return string(nv)
}

var (
	queue         []string
	previousState os.FileInfo
	i             int
	barchartData  []float64
	bcLabels      []string
	statistics    map[string]int
)

func parseLine(b []byte) {
	//127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 123
	re := regexp.MustCompile(`/([a-z]+)`)
	result := re.FindAllString(string(b), -1)
	for i := 0; i < len(result); i++ {
		statistics[result[i]]++
	}
	queue = append(queue, result...)

	barchartData = nil
	bcLabels = nil
	sum := 0
	for _, v := range statistics {
		sum += v
	}
	for k, v := range statistics {
		pourcent := (float64(v) / float64(sum)) * 100
		barchartData = append(barchartData, math.Ceil(pourcent*100)/100)
		bcLabels = append(bcLabels, k)
	}
}

func getLine(initialState os.FileInfo) {
	file, err := os.Open("localfile.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if previousState == nil {
		previousState = initialState
	}
	initialState, _ = os.Stat("localfile.log")

	buf := make([]byte, previousState.Size()+(previousState.Size()-initialState.Size()))

	if previousState.Size() != initialState.Size() && previousState.ModTime() != initialState.ModTime() {
		start := previousState.Size()
		file.ReadAt(buf, start)
		parseLine(buf)
	}
	previousState, _ = os.Stat("localfile.log")

}

func terminalUI() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	initalStat, err := os.Stat("localfile.log")
	if err != nil {
		panic(err)
	}

	statistics = make(map[string]int)

	p := widgets.NewParagraph()
	p.Title = "Text Box"
	p.Text = "PRESS q TO QUIT DEMO"
	p.SetRect(0, 0, 50, 5)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan

	updateParagraph := func(count int) {
		if count%2 == 0 {
			p.TextStyle.Fg = ui.ColorRed
		} else {
			p.TextStyle.Fg = ui.ColorWhite
		}
	}

	listData := []string{
		"[0] gizak/termui",
		"[1] editbox.go",
		"[2] interrupt.go",
		"[3] keyboard.go",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] nsf/termbox-go",
	}

	l := widgets.NewList()
	l.Title = "List"
	l.Rows = listData
	l.SetRect(0, 25, 50, 6)
	l.TextStyle.Fg = ui.ColorYellow

	bc := widgets.NewBarChart()
	bc.Title = "Bar Chart"
	bc.SetRect(50, 0, 80, 25)
	bc.BarWidth = 5
	bc.BarColors[0] = ui.ColorBlue

	draw := func(count int, initalStat os.FileInfo) {
		getLine(initalStat)
		if len(queue) > 0 {
			listData = append(listData, queue[0])
			queue = queue[1:]
		}
		l.Rows = listData[len(listData)-5:]
		bc.Data = barchartData
		bc.Labels = bcLabels

		ui.Render(p, l, bc)
	}

	tickerCount := 1
	draw(tickerCount, initalStat)
	tickerCount++
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			updateParagraph(tickerCount)
			draw(tickerCount, initalStat)
			tickerCount++
		}
	}
}
func main() {
	terminalUI()
}
