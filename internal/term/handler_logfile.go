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
var resHTTPCode = regexp.MustCompile(`(?m)HTTP\/[0-9\.]+\s([0-9]{3})`)

type line struct {
	date     string
	section  string
	httpCode string
}

// LogData manage logfile details
type LogData struct {
	queue      []line
	logfile    string
	recapUsage map[string]int
	dataHandle float64
}

// ParseWithNotify ParseWithNotify
func (l *LogData) ParseWithNotify(errC chan error) error {
	file, _ := os.Open(l.logfile)
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	errWatcher := watcher.Add(l.logfile)
	if errWatcher != nil {
		errC <- errWatcher
	}

	file.Seek(0, os.SEEK_END)
	r := bufio.NewReader(file)

	for {
		by, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			errC <- err
		}

		if err != io.EOF {
			l.parseLine(by)
			l.dataHandle += float64(len(by)) / 1000

			continue
		}
		if err = waitForChange(watcher); err != nil {
			errC <- err
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
func (l *LogData) parseLine(b []byte) {
	sections := reSec.FindAllString(string(b), -1)
	dates := reDate.FindAllStringSubmatch(string(b), -1)
	httpCode := resHTTPCode.FindAllStringSubmatch(string(b), -1)

	queueInfo := make([]line, len(sections))
	/**
		brainstorm about recapUsage behavior.
	**/
	for i := 0; i < len(sections); i++ {
		queueInfo[i].section = sections[i]
		queueInfo[i].date = dates[i][1]
		queueInfo[i].httpCode = httpCode[i][1]
		l.recapUsage[httpCode[i][1]]++
	}

	l.queue = append(l.queue, queueInfo...)
}
