package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ruanda/thinq2mqtt/internal/thinq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		ClientID:     viper.GetString("client_id"),
		ClientSecret: viper.GetString("client_secret"),
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

	fmt.Print("Then paste the URL where the browser is redirected: ")

	var callbackURL string
	fmt.Scanln(&callbackURL)

	res, err := c.Auth.ParseOAuthCallback(callbackURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Refresh token: %s\nAccess token: %s\n", res.RefreshToken, res.AccessToekn)
}
