package cmd

import (
	"emuSync/es"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	emuSync.AddCommand(listDevices)
}

var listDevices = &cobra.Command{
	Use:   "devices",
	Short: "List all devices found by ADB",
	Run: func(cmd *cobra.Command, args []string) {
		es := es.EmuSync{}
		d, err := es.ListDevices()

		if err != nil {
			fmt.Println(err)
			return
		}

		if len(d) == 0 {
			fmt.Println("No devices found")
			return
		}

		var rows []table.Row

		for _, device := range d {
			rows = append(rows, table.Row{device.ID, device.Model})
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "Model"})
		t.AppendRows(rows)
		t.Render()
	},
}
