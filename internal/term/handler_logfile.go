package term

import (
	"os"
	"regexp"
)

// Parse logfile file to get new line added
func (t *Term) Parse(initialState os.FileInfo) {
	file, err := os.Open(t.logfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// no have previous
	// so previous = currentState
	if t.previousState == nil {
		t.previousState = initialState
	}
	initialState, _ = os.Stat(t.logfile)
	// use the file size before a trafic in log file it's observed
	// and substract it to current file size to get a len and cap value of new lines added
	buf := make([]byte,
		t.previousState.Size()+(t.previousState.Size()-initialState.Size()),
		t.previousState.Size()+(t.previousState.Size()-initialState.Size()),
	)
	// if the init size and the current size are different and the update date are different
	if t.previousState.Size() != initialState.Size() && t.previousState.ModTime() != initialState.ModTime() {
		start := t.previousState.Size()
		file.ReadAt(buf, start)
		t.parseLine(buf)
	}
	// new lines are managed, move on tail of file
	t.previousState, _ = os.Stat(t.logfile)
}

// get all section found in new lines added
func (t *Term) parseLine(b []byte) {
	reDate := regexp.MustCompile(`(?m)\[(.+)\]`)
	reSec := regexp.MustCompile(`\s/([a-z]+)`)

	sections := reSec.FindAllString(string(b), -1)
	dates := reDate.FindAllString(string(b), -1)

	queueInfo := make([]line, len(sections))
	for i := 0; i < len(sections); i++ {
		queueInfo[i].section = sections[i]
		queueInfo[i].date = dates[i]
	}

	t.queue = append(t.queue, queueInfo...)
}
