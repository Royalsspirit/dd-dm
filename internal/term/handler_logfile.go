package term

import (
	"os"
	"regexp"
)

// Parse blbalb
func (t *Term) Parse(initialState os.FileInfo) {
	file, err := os.Open(t.logfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if t.previousState == nil {
		t.previousState = initialState
	}
	initialState, _ = os.Stat(t.logfile)

	buf := make([]byte, t.previousState.Size()+(t.previousState.Size()-initialState.Size()))

	if t.previousState.Size() != initialState.Size() && t.previousState.ModTime() != initialState.ModTime() {
		start := t.previousState.Size()
		file.ReadAt(buf, start)
		t.parseLine(buf)
	}

	t.previousState, _ = os.Stat(t.logfile)
}

func (t *Term) parseLine(b []byte) {
	//127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 123
	re := regexp.MustCompile(`/([a-z]+)`)
	result := re.FindAllString(string(b), -1)

	t.sinData = append(t.sinData, float64(len(result)))
	t.queue = append(t.queue, result...)

}
