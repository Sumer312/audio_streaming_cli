package trackcontrols

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	socketcontrols "github.com/sumer312/auditerm/socketControls"
)

var colors []tcell.Color = []tcell.Color{tcell.ColorCadetBlue, tcell.ColorRoyalBlue, tcell.ColorAliceBlue, tcell.ColorCornflowerBlue, tcell.ColorDodgerBlue, tcell.ColorPowderBlue, tcell.ColorMidnightBlue}

func QuitTrack() {
	socket_conection := *socketcontrols.DialConnection()
	defer socket_conection.Close()
	_, err := socket_conection.Write([]byte(`{ "command": ["quit"] }` + "\n"))
	if err != nil {
		log.Panic("write error:", err)
	}
	os.Remove("/tmp/mpvsocket")
}

func TogglePause() {
	socket_conection := *socketcontrols.DialConnection()
	defer socket_conection.Close()
	_, err := socket_conection.Write([]byte(`{ "command": ["get_property", "pause"] }` + "\n"))
	if err != nil {
		log.Panic("write error:", err)
	}
	buffer := make([]byte, 2048)
	n, err := socket_conection.Read(buffer)
	if err != nil {
		log.Panic("read error:", err)
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
		log.Panic("write error:", err)
	}
}

func StartFlicker(ctx context.Context, app *tview.Application, table *tview.Table, rowIndex int, columnCount int) {
	ticker := time.NewTicker(150 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			app.QueueUpdate(func() {
				for i := range columnCount {
					table.GetCell(rowIndex, i).SetTextColor(tcell.ColorYellow)
				}
			})
			return
		case <-ticker.C:
			app.QueueUpdate(func() {
				idx := rand.Intn(len(colors))
				for i := range columnCount {
					table.GetCell(rowIndex, i).SetTextColor(colors[idx])
				}
			})
		}
	}
}
func PlayCurrentTrack(ctx context.Context, app *tview.Application, table *tview.Table, rowIndex int, columnCount int, url string) {
	socketcontrols.MpvInit(url)
	for {
		select {
		case <-ctx.Done():
			app.QueueUpdate(func() {
				for i := range columnCount {
					table.GetCell(rowIndex, i).SetTextColor(tcell.ColorYellow)
				}
			})
			return
		default:
			app.QueueUpdate(func() {
				idx := rand.Intn(len(colors))
				for i := range columnCount {
					table.GetCell(rowIndex, i).SetTextColor(colors[idx])
				}
			})
		}
	}
}
