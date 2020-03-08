package term

import (
	"log"
	"math"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	queue         []string
	previousState os.FileInfo
	i             int
	barchartData  []float64
	bcLabels      []string
	statistics    map[string]int
)

// Term blabla
type Term struct {
	previousState os.FileInfo
	queue         []string
	barchartData  []float64
	bcLabels      []string
	statistics    map[string]int
	sum           int
	logfile		  string
}

// Conf blablba
type Conf struct {
	Logfile string
}
// NewTerm blablab
func NewTerm(conf *Conf) *Term {
	return &Term{
		logfile: conf.Logfile,
	}
}
// Run run
func (t *Term) Run() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	initalStat, err := os.Stat("localfile.log")
	if err != nil {
		panic(err)
	}

	t.statistics = make(map[string]int)

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

	barchartData = nil
	bcLabels = nil
	t.sum = 0

	draw := func(count int, initalStat os.FileInfo) {
		t.Parse(initalStat)
		if len(t.queue) > 0 {
			t.statistics[t.queue[0]]++
			t.sum++
			match := false
			pourcent := (float64(t.statistics[t.queue[0]]) / float64(t.sum)) * 100
			for k, v := range t.bcLabels {
				if v == t.queue[0] {
					t.barchartData[k] = math.Ceil(pourcent*100) / 100
					match = true
				} else {
					pourcent := (float64(t.statistics[v]) / float64(t.sum)) * 100
					t.barchartData[k] = math.Ceil(pourcent*100) / 100
				}
			}
			if !match {
				t.bcLabels = append(t.bcLabels, t.queue[0])
				t.barchartData = append(t.barchartData, math.Ceil(pourcent*100)/100)
			}

			t.queue = t.queue[1:]
			if len(t.queue) != 0 {
				listData = append(listData, t.queue[0])
			}

		}

		l.Rows = listData[len(listData)-5:]
		bc.Data = t.barchartData
		bc.Labels = t.bcLabels

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
