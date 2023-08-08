package cmd

import (
	"fmt"
	"github.com/lightningnetwork/lnd/aezeed"
	"github.com/mustwork/aezeed-recover/internal/ui"
	"github.com/spf13/cobra"
	"os"
)

var (
	devMode  bool
	consent  bool
	InfoText string
)

var rootCmd = &cobra.Command{
	Use:   "aezeed-recover",
	Short: "recover lost or wrong words in aezeed seed phrases",
	Run: func(cmd *cobra.Command, args []string) {
		theme := ui.UniversalTheme()
		config := &ui.Config{
			DevMode:         devMode,
			Consent:         consent,
			Wordlist:        aezeed.DefaultWordList,
			DefaultPassword: "", // 'aezeed' and the empty string are synonymous
			MnemonicLength:  24,
			// cracking more than three words is not feasible
			MnemonicMinLength: 21,
		}
		app := ui.CreateApplication(theme, config, InfoText)
		if err := app.Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&devMode, "dev", "d", false, "enable dev mode (additional menu options)")
	rootCmd.PersistentFlags().BoolVar(&consent, "consent", false, "consent to terms of use (skip info page at startup)")
}

func Execute(s string) {
	InfoText = s
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
