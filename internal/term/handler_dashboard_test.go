package term

import (
	"testing"
	"time"

	"github.com/gizak/termui/v3/widgets"
	"github.com/stretchr/testify/assert"
)

// TestAlert test dashboard with a threshold and a custom start date
func TestAlert(t *testing.T) {
	te := Term{
		start:     time.Now().Add(-120 * time.Second),
		threshold: 10,
	}
	d := dashboard{
		pa: widgets.NewParagraph(),
	}
	drawAlert(&te, &d)
	assert.Contains(t, d.pa.Text, "alert", "should contain alert")
}

// TestBartChart test a barchart percentage
func TestBartChart(t *testing.T) {
	te := Term{
		statistics: map[string]int{"/toto": 0, "/tata": 0},
		sum:        1,
	}

	drawBarchart(&te)
	assert.Equal(t, te.barchartData[0], float64(100), "should equal to 100")
	te.sum = 2
	drawBarchart(&te)

	assert.Equal(t, te.barchartData[0], float64(50), "should equal to 50")

}
