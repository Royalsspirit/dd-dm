package term

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Term blabla
type Term struct {
	barchartData []float64
	bcLabels     []string
	statistics   map[string]int
	sum          int
	logConf      LogData
	sinData      []float64
	threshold    int
	start        time.Time
}

type dashboard struct {
	p  *widgets.Plot
	b  *widgets.BarChart
	l  *widgets.List
	pa *widgets.Paragraph
	p2 *widgets.Paragraph
}

// Conf contain a cli parameter required to run a dashboard
type Conf struct {
	Logfile   string
	Threshold int
}

var orderHTTPCode = []string{
	"5xx",
	"4xx",
	"3xx",
	"2xx",
	"1xx",
}
var tplHTTPUsade = map[string]int{
	"5xx": 0,
	"4xx": 0,
	"3xx": 0,
	"2xx": 0,
	"1xx": 0,
}

// NewTerm create a newTerm configuration
func NewTerm(conf *Conf) *Term {
	return &Term{
		logConf: LogData{
			logfile: conf.Logfile,
		},
		threshold: conf.Threshold,
		start:     time.Now(),
	}
}

func (t *Term) makeDashboard() *dashboard {
	p := widgets.NewParagraph()
	p.Title = "HTTP Usage"
	p.SetRect(0, 0, 80, 5)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan

	p2 := widgets.NewParagraph()
	p2.Title = "History"
	p2.SetRect(80, 0, 119, 30)
	p2.TextStyle.Fg = ui.ColorWhite
	p2.BorderStyle.Fg = ui.ColorCyan

	// default list value
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
	l.Title = "Log"
	l.Rows = listData
	l.SetRect(0, 5, 50, 20)
	l.TextStyle.Fg = ui.ColorYellow

	bc := widgets.NewBarChart()
	bc.Title = "Pourcentage of each section"
	bc.SetRect(50, 5, 80, 20)
	bc.BarWidth = 6
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	// make a slice with a len of 73 corresponding to barchart width
	t.sinData = make([]float64, 73)

	lc2 := widgets.NewPlot()
	lc2.Title = "Number of requests"
	lc2.Data = make([][]float64, 1)
	lc2.Data[0] = t.sinData
	lc2.SetRect(0, 20, 80, 30)
	lc2.AxesColor = ui.ColorWhite
	lc2.LineColors[0] = ui.ColorYellow

	result := dashboard{b: bc, p: lc2, l: l, pa: p, p2: p2}
	return &result
}

func drawBarchart(t *Term) {
	t.statistics[t.logConf.queue[0].section]++
	match := false
	pourcent := (float64(t.statistics[t.logConf.queue[0].section]) / float64(t.sum)) * 100
	for k, v := range t.bcLabels {
		if v == t.logConf.queue[0].section {
			t.barchartData[k] = math.Ceil(pourcent*100) / 100
			match = true
		} else {
			pourcent := (float64(t.statistics[v]) / float64(t.sum)) * 100
			t.barchartData[k] = math.Ceil(pourcent*100) / 100
		}
	}
	if !match {
		t.bcLabels = append(t.bcLabels, t.logConf.queue[0].section)
		t.barchartData = append(t.barchartData, math.Ceil(pourcent*100)/100)
	}
}

func drawAlert(t *Term, d *dashboard) {
	currentTime := time.Now()
	d.p2.Text += t.logConf.alert.declare(currentTime, t.threshold)
}

func drawList(d *dashboard, t *Term) {
	if len(t.logConf.queue) != 0 {
		date, _ := time.Parse("02/Jan/2006:15:04:05 -0700", t.logConf.queue[0].date)
		d.l.Rows = append(d.l.Rows, fmt.Sprint(t.logConf.queue[0].httpCode, " - ", date.Local(), " - ", t.logConf.queue[0].section))
		d.l.Rows = d.l.Rows[1:]
	}
}

func drawLine(d *dashboard, t *Term) {
	d.p.Data[0] = t.sinData
	if len(t.sinData) > 72 {
		t.sinData = t.sinData[2:]
	} else {
		t.sinData = append(t.sinData, t.sinData[len(t.sinData)-1])
	}
}

// every 10 seconds draw the dashboard
func drawDashboard(t *Term, d *dashboard) {
	if len(t.logConf.queue) > 0 {
		for len(t.logConf.queue) > 0 {
			t.sum++

			drawBarchart(t)

			drawList(d, t)

			t.logConf.queue = t.logConf.queue[1:]
		}
	}

	drawAlert(t, d)

	t.sinData = append(t.sinData, float64(t.sum))

	drawLine(d, t)
	d.b.Data = t.barchartData
	d.b.Labels = t.bcLabels

	var httpCodeDetails string
	// try to imprive this code
	for k, v := range orderHTTPCode {
		var currentValue string
		if t.logConf.recapUsage[v] != 0 {
			currentValue = strconv.Itoa(t.logConf.recapUsage[v])
			// display value only once
			t.logConf.recapUsage[v] = 0
		} else {
			currentValue = strconv.Itoa(tplHTTPUsade[v])
		}

		httpCodeDetails += v + ": " + currentValue + "     "

		if k > 0 && k%4 == 0 {
			httpCodeDetails += "\n"
		}
	}

	t.logConf.totalDataHandle += t.logConf.dataHandle

	d.pa.Text = fmt.Sprint(httpCodeDetails, "Total Requests RX: ", t.logConf.totalDataHandle, " B", "              Rx/s: ", t.logConf.dataHandle)

	ui.Render(d.p, d.l, d.b, d.pa)

	t.logConf.dataHandle = 0

}

// Run dashboard
func (t *Term) Run() error {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	// create a chan to manage error in go routine
	errc := make(chan error)

	go t.logConf.ParseWithNotify(errc, t.threshold)

	t.statistics = make(map[string]int)
	t.logConf.recapUsage = tplHTTPUsade
	dashboard := t.makeDashboard()

	updateParagraph := func() {
		dashboard.p2.TextStyle.Fg = ui.ColorWhite
		ui.Render(dashboard.p2)
	}

	t.start = time.Now()
	t.sum = 0
	// init dashboard
	drawDashboard(t, dashboard)
	uiEvents := ui.PollEvents()

	tickerUI := time.NewTicker(time.Second * 10).C
	tickerParser := time.NewTicker(time.Second).C

	for {
		select {
		case err := <-errc:
			return err
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			}
		case <-tickerUI:
			drawDashboard(t, dashboard)
		case <-tickerParser:
			updateParagraph()
		}
	}
}
