package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	trackcontrols "github.com/sumer312/auditerm/trackControls"
)

func main() {
	logFile := SetupLogger()
	defer logFile.Close()
	log.Print("App started")
	app := tview.NewApplication()
	table := tview.NewTable()
	table.SetBackgroundColor(tcell.ColorNone)
	var songs []string = []string{"1,The pointer sisters,Hot-together,4:14,https://www.youtube.com/watch?v=7k0eEdoZ9JI", "2,Aimer,Kiro,6:51,https://www.youtube.com/watch?v=M3J1KRD1H1Q", "3,Hige,Mixed Nuts,3:32,https://www.youtube.com/watch?v=Wf8XuAoCjN8"}

	ctx, cancel := context.WithCancel(context.Background())

	for i := range songs {
		song := strings.Split(songs[i], ",")
		for c := range len(song) {
			cell := tview.NewTableCell(song[c]).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter)
			table.SetCell(i, c, cell)
		}
	}

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		if ctx.Value("row") == nil || ctx.Value("row").(int) != row {
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "row", row)
			url := table.GetCell(row, table.GetColumnCount()-1).Text
			go trackcontrols.PlayCurrentTrack(ctx, app, table, row, table.GetColumnCount(), url)
		} else {
			go trackcontrols.TogglePause()
		}
	})
	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		f, err := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		log.SetOutput(f)
		panic(err)
	}
}
