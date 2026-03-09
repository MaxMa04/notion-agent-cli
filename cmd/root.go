package cmd

import (
	"fmt"
	"os"

	"github.com/MaxMa04/notion-agent-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	debugMode    bool
	profileFlag  string
	// Version is set by goreleaser ldflags
	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "notion-agent",
	Short: "Notion Agent CLI — work seamlessly with Notion from the command line",
	Long: `Notion Agent CLI — work seamlessly with Notion from the command line.

Manage pages, databases, blocks, and more without leaving your terminal.
Built for developers and AI agents with native multi-profile support.

Use --profile to select a specific auth profile per command,
or set NOTION_PROFILE to auto-select by environment.`,
	Version:       Version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "", "Output format: json, md, table, text (default: auto)")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Show HTTP request/response details")
	rootCmd.PersistentFlags().StringVar(&profileFlag, "profile", "", "Use a specific auth profile for this command")

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(pageCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(blockCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(commentCmd)
	rootCmd.AddCommand(fileCmd)
}

// getToken returns the Notion API token from env, profile flag, or config file.
// Priority: NOTION_TOKEN env > --profile flag > NOTION_PROFILE env > current_profile in config.
func getToken() (string, error) {
	// 1. Direct token via environment variable
	if token := os.Getenv("NOTION_TOKEN"); token != "" {
		return token, nil
	}

	// 2-4. Profile-based resolution from config
	cfg, err := config.Load()
	if err == nil {
		// 2. --profile CLI flag
		if profileFlag != "" {
			if p := cfg.GetProfile(profileFlag); p != nil && p.Token != "" {
				return p.Token, nil
			}
			return "", fmt.Errorf("profile %q not found. Run 'notion-agent auth login --profile %s' first", profileFlag, profileFlag)
		}

		// 3. NOTION_PROFILE environment variable
		if envProfile := os.Getenv("NOTION_PROFILE"); envProfile != "" {
			if p := cfg.GetProfile(envProfile); p != nil && p.Token != "" {
				return p.Token, nil
			}
			return "", fmt.Errorf("profile %q not found. Run 'notion-agent auth login --profile %s' first", envProfile, envProfile)
		}

		// 4. Current profile from config
		if p := cfg.GetCurrentProfile(); p != nil && p.Token != "" {
			return p.Token, nil
		}
	}

	return "", fmt.Errorf("not authenticated. Run 'notion-agent auth login' or set NOTION_TOKEN")
}
