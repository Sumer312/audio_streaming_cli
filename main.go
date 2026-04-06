package main

import (
	"context"
	"encoding/json"
	"log"
	/* "math/rand" */
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var colors []tcell.Color = []tcell.Color{tcell.ColorCadetBlue, tcell.ColorRoyalBlue, tcell.ColorAliceBlue, tcell.ColorCornflowerBlue, tcell.ColorDodgerBlue, tcell.ColorPowderBlue, tcell.ColorMidnightBlue}

func dialConnection() *net.Conn {
	socket_conection, err := net.Dial("unix", "/tmp/mpvsocket")
	if err != nil {
		return nil
	}
	return &socket_conection
}

func mpvInit(url string) *os.Process {
	cmd := exec.Command("mpv", "--input-ipc-server=/tmp/mpvsocket", "--no-video", url)
	cmd.Start()
	proc, err := os.FindProcess(cmd.Process.Pid)
	if err != nil {
		panic(err.Error())
	}
	return proc
}

func quitTrack() {
	socket_conection := *dialConnection()
	defer socket_conection.Close()
	_, err := socket_conection.Write([]byte(`{ "command": ["quit"] }` + "\n"))
	if err != nil {
		log.Fatal("write error:", err)
	}
}

func togglePause() {
	socket_conection := *dialConnection()
	defer socket_conection.Close()
	_, err := socket_conection.Write([]byte(`{ "command": ["get_property", "pause"] }` + "\n"))
	if err != nil {
		log.Fatal("write error:", err)
	}
	buffer := make([]byte, 2048)
	n, err := socket_conection.Read(buffer)
	if err != nil {
		log.Fatal("read error:", err)
	}
	tmp := buffer[:n]

	var buffer_ouput map[string]any
	json.Unmarshal(tmp, &buffer_ouput)

	if buffer_ouput["data"].(bool) {
		_, err = socket_conection.Write([]byte(`{ "command": ["set_property", "pause", false] }` + "\n"))
	} else {
		_, err = socket_conection.Write([]byte(`{ "command": ["set_property", "pause", true] }` + "\n"))
	}
	if err != nil {
		log.Fatal("write error:", err)
	}
}

func playCurrentTrack(ctx context.Context, app *tview.Application, table *tview.Table, row int, columnCount int, url string) {
	if dialConnection() != nil {
		quitTrack()
	}
	mpvInit(url)
	/* for { */
	/* 	app.QueueUpdate(func() { */
	/* 		select { */
	/* 		case <-ctx.Done(): */
	/* 			quitTrack() */
	/* 			for i := range columnCount { */
	/* 				table.GetCell(row, i).SetTextColor(tcell.ColorYellow) */
	/* 			} */
	/* 			return */
	/* 		default: */
	/* 			idx := rand.Intn(len(colors)) */
	/* 			for i := range columnCount { */
	/* 				table.GetCell(row, i).SetTextColor(colors[idx]) */
	/* 			} */
	/* 		} */
	/* 	}) */
	/* } */
}

func main() {
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
			cell := table.GetCell(row, table.GetColumnCount()-1)
			go playCurrentTrack(ctx, app, table, row, table.GetColumnCount(), cell.Text)
		} else {
			go togglePause()
		}
	})
	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
