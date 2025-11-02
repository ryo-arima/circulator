package client

import (
	"github.com/ryo-arima/circulator/pkg/client/controller"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/spf13/cobra"
)

type BaseCmd struct {
	Bootstrap *cobra.Command
	Create    *cobra.Command
	Get       *cobra.Command
	Update    *cobra.Command
	Delete    *cobra.Command
	Config    config.BaseConfig // Dependency injection for config
}

// InitRootCmd creates the root command with config dependency injection
func InitRootCmd(baseConfig config.BaseConfig) *cobra.Command {
	var output string
	rootCmd := &cobra.Command{
		Use:   "circulator",
		Short: "'circulator' is a CLI tool to manage circulator resources",
		Long:  `''`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			controller.SetOutputFormat(output)
			// Log command execution
			baseConfig.Logger.INFO(config.CBCE, "Command executed", map[string]interface{}{
				"command": cmd.Name(),
				"args":    args,
				"output":  output,
			})
		},
	}
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "Output format: table|json|yaml")
	return rootCmd
}

// InitBaseCmd creates base commands with config dependency injection
func InitBaseCmd(baseConfig config.BaseConfig) BaseCmd {
	baseConfig.Logger.DEBUG(config.CBIBC, "Initializing base commands")

	bootstrapCmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap resources",
		Long:  `Bootstrap resources for circulator`,
		Run: func(cmd *cobra.Command, args []string) {
			baseConfig.Logger.INFO(config.CBBCC, "Bootstrap command called")
			cmd.Help()
		},
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create resources",
		Long:  `Create resources for circulator`,
		Run: func(cmd *cobra.Command, args []string) {
			baseConfig.Logger.INFO(config.CBCCC, "Create command called")
			cmd.Help()
		},
	}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get resources",
		Long:  `Get resources from circulator`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update resources",
		Long:  `Update resources in circulator`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
		Long:  `Delete resources from circulator`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	baseCmd := BaseCmd{
		Bootstrap: bootstrapCmd,
		Create:    createCmd,
		Get:       getCmd,
		Update:    updateCmd,
		Delete:    deleteCmd,
		Config:    baseConfig, // Store config in BaseCmd
	}
	return baseCmd
}

func Client(conf config.BaseConfig) {
	conf.Logger.INFO(config.CBSCA, "Starting client application")

	rootCmd := InitRootCmd(conf)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	baseCmd := InitBaseCmd(conf)

	// Add base commands to root
	rootCmd.AddCommand(baseCmd.Bootstrap)
	rootCmd.AddCommand(baseCmd.Create)
	rootCmd.AddCommand(baseCmd.Get)
	rootCmd.AddCommand(baseCmd.Update)
	rootCmd.AddCommand(baseCmd.Delete)

	conf.Logger.DEBUG(config.CBACR, "All commands registered")
	rootCmd.Execute()
}
