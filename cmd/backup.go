package cmd

import (
	"emuSync/es"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	backup.Flags().StringP("device", "d", "", "Device Serial ID")
	backup.Flags().BoolP("backup-roms", "r", false, "Should ROMs be backed up?")

	_ = backup.MarkFlagRequired("device")

	emuSync.AddCommand(backup)
}

var backup = &cobra.Command{
	Use:   "backup",
	Short: "Backs up all configured directories",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("device")
		backupROMs, _ := cmd.Flags().GetBool("backup-roms")

		es := es.EmuSync{}

		d, err := es.GetDeviceWithConfig(id)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = es.BackupDevice(d, backupROMs)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("device")
		es := es.EmuSync{}

		if !es.DoesConfigExist(id) {
			return fmt.Errorf("config file does not exist for device: %s", id)
		}

		return nil
	},
}
