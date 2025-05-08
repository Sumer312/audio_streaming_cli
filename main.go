package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var colors []tcell.Color = []tcell.Color{tcell.ColorPink, tcell.ColorLightCoral, tcell.ColorIvory, tcell.ColorRed, tcell.ColorSpringGreen, tcell.ColorLightCyan, tcell.ColorRoyalBlue}

func mpvInit(url string) *os.Process {
	cmd := exec.Command("mpv", "--input-ipc-server=/tmp/mpvsocket", "--no-video", url)
	cmd.Start()
	proc, err := os.FindProcess(cmd.Process.Pid)
	if err != nil {
		panic(err.Error())
	}
	return proc
}
func killTrack(pid *os.Process) bool {
	err := pid.Kill()
	if err != nil {
		panic(err.Error())
	}
	return err == nil
}
func pauseTrack(shouldPause bool) {
	socat := exec.Command("socat", "-", "tmp/mpvsocket")
	stdin, err := socat.StdinPipe()
	if err != nil {
		fmt.Println(err.Error())
	}
	if shouldPause {
		io.WriteString(stdin, `{ "command": ["set_property", "pause", true] }`)
	} else {
		io.WriteString(stdin, `{ "command": ["set_property", "pause", false] }`)
	}
	stdin.Close()
	socat.Wait()
}
func playCurrentTrack(ctx context.Context, app *tview.Application, table *tview.Table, pauseChan chan bool, row int, columnCount int, url string) {
	var pid os.Process = *mpvInit(url)
	for {
		app.QueueUpdate(func() {
			select {
			case <-ctx.Done():
				killTrack(&pid)
				for i := range columnCount {
					table.GetCell(row, i).SetTextColor(tcell.ColorYellow)
				}
				return
			case flag := <-pauseChan:
				if flag {
					pauseTrack(true)
				} else {
					pauseTrack(false)
				}
			default:
				idx := rand.Intn(len(colors))
				for i := range columnCount {
					table.GetCell(row, i).SetTextColor(colors[idx])
				}
			}
		})
	}
}

func main() {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(true)
	table.SetBackgroundColor(tcell.ColorNone)
	song1 := strings.Split("1,The pointer sisters,Hot-together,4:14,https://www.youtube.com/watch?v=7k0eEdoZ9JI", ",")
	song2 := strings.Split("2,Aimer,Kiro,6:51,https://www.youtube.com/watch?v=M3J1KRD1H1Q", ",")
	song3 := strings.Split("3,Hige,Mixed Nuts,3:32,https://www.youtube.com/watch?v=Wf8XuAoCjN8", ",")

	cols := len(song1)
	ctx := context.Background()
	callCtx, cancel := context.WithCancel(ctx)
	pauseChan := make(chan bool, 1)
	trackMap := [2]int{-1, -1}
	for c := range cols {
		cell := tview.NewTableCell(song1[c]).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter)
		table.SetCell(0, c, cell)
	}
	for c := range cols {
		cell := tview.NewTableCell(song2[c]).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter)
		table.SetCell(1, c, cell)
	}
	for c := range cols {
		cell := tview.NewTableCell(song3[c]).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter)
		table.SetCell(2, c, cell)
	}
	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		if callCtx.Value("row") == nil || callCtx.Value("row").(int) != row {
			cancel()
			trackMap[0] = row
			trackMap[1] = 1
			pauseChan = make(chan bool, 1)
			callCtx, cancel = context.WithCancel(ctx)
			callCtx = context.WithValue(callCtx, "row", row)
			cell := table.GetCell(row, table.GetColumnCount()-1)
			go playCurrentTrack(callCtx, app, table, pauseChan, row, table.GetColumnCount(), cell.Text)
		} else {
			panic("hhiiiii")
			/* if trackMap[0] != -1 { */
			/* 	if trackMap[1] == 0 { */
			/* 		trackMap[1] = 1 */
			/* 		pauseChan <- false */
			/* 	} else { */
			/* 		trackMap[1] = 0 */
			/* 		pauseChan <- true */
			/* 	} */
			/* } */
		}
	})
	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
