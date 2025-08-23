package cmd

import (
	"emuSync/es"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	listFiles.Flags().StringP("device", "d", "", "Device ID to target")

	emuSync.AddCommand(listFiles)
}

var listFiles = &cobra.Command{
	Use:   "ls <path>",
	Short: "List all files for a given path",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		es := es.EmuSync{}

		id, err := cmd.Flags().GetString("device")
		if err != nil {
			fmt.Println(err)
			return
		}

		path := args[0]
		d, err := es.GetDevice(id)

		if err != nil {
			fmt.Println(err)
			return
		}

		files, err := es.ListFiles(d, path)
		if err != nil {
			fmt.Println(err)
			return
		}

		var rows []table.Row

		for _, file := range files {
			rows = append(rows, table.Row{file.Name, file.Path, file.Size, file.LastModified})
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "Path", "Size", "Last Modified"})
		t.AppendRows(rows)
		t.Render()
	},
}
