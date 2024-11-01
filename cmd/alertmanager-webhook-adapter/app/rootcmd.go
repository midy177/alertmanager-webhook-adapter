package app

import (
	"flag"
	"fmt"
	"os"

	"alertmanager-webhook-adapter/cmd/alertmanager-webhook-adapter/app/options"
	"alertmanager-webhook-adapter/pkg/version"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	o := options.NewAppOptions()

	rootCmd := &cobra.Command{
		Use:   "alertmanager-webhook-adapter",
		Short: "alertmanager-webhook-adapter",
		Long:  `alertmanager-webhook-adapter`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if o.Version {
				fmt.Printf("Version: %s\n", version.Version)
				fmt.Printf("Commit: %s\n", version.Commit)
				fmt.Printf("BuildAt: %s\n", version.BuildAt)
				return nil
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			if o.Version {
				return
			}

			if err := o.Run(); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
				return
			}
		},
	}

	rootCmd.Flags().StringVarP(&o.Addr, "listen-address", "l", "0.0.0.0:8090", "the address to listen")
	rootCmd.Flags().StringVarP(&o.Signature, "signature", "s", "未知", "the signature")
	rootCmd.Flags().StringVarP(&o.TmplDir, "tmpl-dir", "d", "", "the tmpl dir")
	rootCmd.Flags().StringVarP(&o.TmplName, "tmpl-name", "t", "", "the tmpl name")
	rootCmd.Flags().StringVarP(&o.TmplDefault, "tmpl-default", "n", "", "the default tmpl name")
	rootCmd.Flags().StringVarP(&o.TmplLang, "tmpl-lang", "", "", "the language for template filename")
	rootCmd.Flags().BoolVarP(&o.Version, "version", "v", false, "show version")
	rootCmd.Flags().BoolVarP(&o.Debug, "debug", "", false, "enable verbose output ")

	rootCmd.Flags().AddGoFlagSet(flag.CommandLine)

	return rootCmd
}
