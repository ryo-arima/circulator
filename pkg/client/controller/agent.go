package controller

import (
	"github.com/ryo-arima/circulator/pkg/client/usecase"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/request"
	"github.com/spf13/cobra"
)

func InitAgentCmd(conf config.BaseConfig) *cobra.Command {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Agent operations",
		Long:  "Manage agents through the API",
	}

	// Initialize usecase
	agentUsecase := usecase.NewAgentUsecase(conf)

	// Add subcommands
	agentCmd.AddCommand(bootstrapAgentCmd(agentUsecase))
	agentCmd.AddCommand(getAgentCmd(agentUsecase))
	agentCmd.AddCommand(createAgentCmd(agentUsecase))
	agentCmd.AddCommand(updateAgentCmd(agentUsecase))
	agentCmd.AddCommand(deleteAgentCmd(agentUsecase))

	return agentCmd
}

func bootstrapAgentCmd(agentUsecase usecase.AgentUsecase) *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap agent data",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := request.AgentRequest{}
			result := agentUsecase.Bootstrap(req, format)
			cmd.Printf("%s\n", result)
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "json", "Output format (json, yaml, table)")

	return cmd
}

func getAgentCmd(agentUsecase usecase.AgentUsecase) *cobra.Command {
	var format string
	var uuid string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get agent data",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := request.AgentRequest{UUID: uuid}
			result := agentUsecase.Get(req, format)
			cmd.Printf("%s\n", result)
			return nil
		},
	}
	cmd.Flags().StringVarP(&uuid, "uuid", "u", "", "UUID filter (optional)")
	cmd.Flags().StringVar(&format, "format", "json", "Output format (json, yaml, table)")

	return cmd
}

func createAgentCmd(agentUsecase usecase.AgentUsecase) *cobra.Command {
	var format string
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := request.AgentRequest{Name: name}
			result := agentUsecase.Create(req, format)
			cmd.Printf("%s\n", result)
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name (optional)")

	cmd.Flags().StringVar(&format, "format", "json", "Output format (json, yaml, table)")

	return cmd
}

func updateAgentCmd(agentUsecase usecase.AgentUsecase) *cobra.Command {
	var format string
	var uuid string
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := request.AgentRequest{UUID: uuid, Name: name}
			result := agentUsecase.Update(req, format)
			cmd.Printf("%s\n", result)
			return nil
		},
	}
	cmd.Flags().StringVarP(&uuid, "uuid", "u", "", "UUID (required)")
	cmd.MarkFlagRequired("uuid")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name (optional)")

	cmd.Flags().StringVar(&format, "format", "json", "Output format (json, yaml, table)")

	return cmd
}

func deleteAgentCmd(agentUsecase usecase.AgentUsecase) *cobra.Command {
	var format string
	var uuid string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := request.AgentRequest{UUID: uuid}
			result := agentUsecase.Delete(req, format)
			cmd.Printf("%s\n", result)
			return nil
		},
	}
	cmd.Flags().StringVarP(&uuid, "uuid", "u", "", "UUID (required)")
	cmd.MarkFlagRequired("uuid")

	cmd.Flags().StringVar(&format, "format", "json", "Output format (json, yaml, table)")

	return cmd
}
