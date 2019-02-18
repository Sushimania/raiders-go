package main

import (
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"github.com/spf13/cobra"
	"raiders-go/service"
	log "github.com/sirupsen/logrus"
)

var (
	version = "0.0.1"
)

func main() {
	// Get EOS Account from local file
	eosAccountName := service.GetAccountFromWallet()

	// Get machine id
	machineId, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}

	var cmdStart = &cobra.Command{
		Use:   "start",
		Short: "Start searching for Bitcoin",
		Long: `Start searching for Bitcoin`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			service.SetGenerate(eosAccountName, machineId)
		},
	}

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Returns the version of this client",
		Long: `Returns the version of this client`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	var cmdImportAccount = &cobra.Command{
		Use:   "import [EOS Account Name]",
		Short: "Adds a EOS Account name",
		Long: `Adds a EOS Account name`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			service.ImportAccount(args[0])
		},
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdStart)
	rootCmd.AddCommand(cmdVersion)
	rootCmd.AddCommand(cmdImportAccount)
	rootCmd.Execute()
}