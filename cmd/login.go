package cmd

import (
	"context"
	"os"
	"fmt"
	"github.com/spf13/viper"
	"github.com/ruanda/thinq2mqtt/internal/thinq"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use: "login",
	Run: loginRun,
}

func loginRun(cmd *cobra.Command, args []string) {
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
	u, err := c.Auth.GetOAuthURL()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Log in here:\n%v\n", u)
}