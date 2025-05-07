package main

import (
	/* "fmt" */
	"context"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var colors []tcell.Color = []tcell.Color{tcell.ColorMaroon, tcell.ColorPurple, tcell.ColorLightCyan, tcell.ColorPink, tcell.ColorLimeGreen}

func playCurrentTrack(ctx context.Context, table *tview.Table, row int, columnCount int) {
	idx := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for i := range columnCount {
				table.GetCell(row, i).SetTextColor(colors[idx%len(colors)])
			}
			time.Sleep(50 * time.Millisecond)
			idx++
		}
	}
}
func pauseCurrentTrack(table *tview.Table, row int, columnCount int) {
	for i := range columnCount {
		table.GetCell(row, i).SetTextColor(tcell.ColorYellow)
	}
}

func main() {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(true)
	table.SetBackgroundColor(tcell.ColorNone)
	song := strings.Split("1 Aimer Re-frain 3:09", " ")
	cols, rows := len(song), 8
	var cancel context.CancelFunc
	ctx := context.Background()
	callCtx, cancel := context.WithCancel(ctx)
	for r := range rows {
		for c := range cols {
			cell := tview.NewTableCell(song[c]).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter)
			table.SetCell(r, c, cell)
		}
	}
	var current_track [2]int = [2]int{0, 0}
	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		if current_track[1] == row {
			if current_track[0] == 0 {
				current_track[0] = 1
			} else {
				current_track[0] = 0
			}
			cancel()
			pauseCurrentTrack(table, row, table.GetColumnCount())
			return
		}
		if current_track[0] != 0 {
			cancel()
			callCtx, cancel = context.WithCancel(ctx)
		}
		current_track[0] = 1
		current_track[1] = row
		go playCurrentTrack(callCtx, table, row, table.GetColumnCount())
	})
	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
