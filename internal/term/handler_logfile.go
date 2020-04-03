package term

import (
	"bufio"
	"io"
	"os"
	"regexp"

	"github.com/fsnotify/fsnotify"
)

var reDate = regexp.MustCompile(`(?m)\[(.+)\]`)
var reSec = regexp.MustCompile(`\s/([a-z]+)`)

// Parse logfile file to get new line added
func (t *Term) Parse(initialState os.FileInfo) {
	// no have previous
	// so previous = currentState
	if t.previousState == nil {
		t.previousState = initialState
	}
	initialState, _ = os.Stat(t.logfile)
	// use the file size before a trafic in log file it's observed
	// and substract it to current file size to get a len and cap value of new lines added
	buf := make([]byte,
		initialState.Size()-t.previousState.Size(),
		initialState.Size()-t.previousState.Size(),
	)
	// if the init size and the current size are different and the update date are different
	if t.previousState.Size() != initialState.Size() && t.previousState.ModTime() != initialState.ModTime() {
		file, err := os.Open(t.logfile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		start := t.previousState.Size()
		file.ReadAt(buf, start)
		t.parseLine(buf)
	}
	// new lines are managed, move on tail of file
	t.previousState, _ = os.Stat(t.logfile)
}

// ParseWithNotify ParseWithNotify
func (t *Term) ParseWithNotify() error {
	file, _ := os.Open(t.logfile)
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()
	_ = watcher.Add(t.logfile)

	file.Seek(0, os.SEEK_END)
	r := bufio.NewReader(file)
	for {
		by, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}
		t.parseLine(by)
		if err != io.EOF {
			continue
		}
		if err = waitForChange(watcher); err != nil {
			return err
		}
	}
}

func waitForChange(w *fsnotify.Watcher) error {
	for {
		select {
		case event := <-w.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				return nil
			}
		case err := <-w.Errors:
			return err
		}
	}
}

// get all section found in new lines added
func (t *Term) parseLine(b []byte) {
	sections := reSec.FindAllString(string(b), -1)
	dates := reDate.FindAllString(string(b), -1)

	queueInfo := make([]line, len(sections))
	for i := 0; i < len(sections); i++ {
		queueInfo[i].section = sections[i]
		queueInfo[i].date = dates[i]
	}

	t.queue = append(t.queue, queueInfo...)
}
