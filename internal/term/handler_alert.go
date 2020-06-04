package term

import (
	"strconv"
	"time"
)

type state struct {
	history     []time.Time
	highTraffic bool
}

func (s *state) declare(currentTime time.Time, threshold int) string {
	var savedIndex int = 0
	var message string
	for k, v := range s.history {
		if currentTime.Sub(v).Seconds() > 120 {
			savedIndex = k
		}
	}

	s.history = s.history[savedIndex:]

	if (float64(len(s.history)) / 120) > float64(threshold) {
		if s.highTraffic == false {
			s.highTraffic = true

			message = "High traffic generated an alert - hits =" + strconv.Itoa(len(s.history)) + ", triggered at " + time.Now().String() + "\n"
		}
	} else {
		if s.highTraffic == true {
			s.highTraffic = false

			message = "Recovered triggered at " + time.Now().String() + "\n"

		}
	}
	return message
}
