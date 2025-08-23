package cmd

import (
	"emuSync/es"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	deviceCmd.AddCommand(initDevice)

	emuSync.AddCommand(deviceCmd)
	emuSync.AddCommand(listDevices)
}

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Device management commands",
	Long:  "Commands for managing device configurations and settings",
}

var initDevice = &cobra.Command{
	Use:   "init <device serial>",
	Short: "Initialize a device for use with emuSync",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		es := es.EmuSync{}
		id := args[0]
		d, err := es.InitDevice(id)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Device %s [%s] initialized successfully", d.ID, d.Model)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		es := es.EmuSync{}

		if es.DoesConfigExist(id) {
			return fmt.Errorf("config file already exists for device: %s", id)
		}

		return nil
	},
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
