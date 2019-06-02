package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ruanda/thinq2mqtt/internal/thinq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test",
	Run:   testRun,
}

func testRun(cmd *cobra.Command, args []string) {
	cfg := thinq.Config{
		CountryCode:  viper.GetString("country"),
		LanguageCode: viper.GetString("language"),
		ServiceCode:  viper.GetString("service"),
		ClientID:     viper.GetString("client"),
	}
	c, err := thinq.NewClient(cfg, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = c.Gateway.Discover(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
