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
