package printer

import (
	"fmt"
	"strconv"
	"time"
)

type printer struct {
	colors *textColors
}

func NewPrinter() IPrinter {
	return &printer{
		colors: NewColors(),
	}
}

func (p *printer) Preview(competitor string) {
	fmt.Println(p.colors.Cyan + competitor + p.colors.Reset)

	t := time.Now()
	fmt.Println(t.Format("2006-01-02 15:04:05"))
}

func (p *printer) Flats(month string, countAllFlats int, newFlats []string) {
	fmt.Println(
		month +
		": " +
		p.colors.Yellow + strconv.Itoa(countAllFlats) + " results" + p.colors.Reset,
	)
	if len(newFlats) > 0 {
		fmt.Println(
			"    Status:", 
			p.colors.Green + strconv.Itoa(len(newFlats)) + " new items found!" + p.colors.Reset,
		)
		for i, flat := range newFlats {
			var color string
			if i % 2 == 1 {
				color = p.colors.Yellow
			} else {
				color = p.colors.Blue
			}
			fmt.Println()
			fmt.Println("        Link:", color + flat + p.colors.Reset)
		}
	}
	fmt.Println()
}