package cmd

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	// juju/persistent-cookiejar

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login user password",
	Short: "Login to leetcode account with Github",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Missing username and/or password")
			return
		}
		username := args[0]
		// password := args[1]

		var client http.Client

		// leetcodeLogin := "https://leetcode.com/accounts/github/login/?next=%2F"
		githubLogin := "https://github.com/login"

		res, err := client.Get(githubLogin)
		if err != nil {
			print(err)
		}

		defer res.Body.Close()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			print(err)
		}
		bodyString := string(bodyBytes)

		authRegex, _ := regexp.Compile(`name="authenticity_token" value="(.*?)"`)
		authToken := authRegex.FindString(bodyString)

		idRegex, _ := regexp.Compile(`name="ga_id" value="(.*?)"`)
		gaId := idRegex.FindString(bodyString)

		timestampRegex, _ := regexp.Compile(`name="timestamp" value="(.*?)"`)
		timestamp := timestampRegex.FindString(bodyString)

		timestampSecretRegex, _ := regexp.Compile(`name="timestamp_secret" value="(.*?)"`)
		timestampSecret := timestampSecretRegex.FindString(bodyString)

		if !(authToken != "" && timestamp != "" && timestampSecret != "") {
			fmt.Println("Couldn't find required fields")
			return
		}

		fmt.Println(authToken)
		fmt.Println(gaId)
		fmt.Println(timestamp)
		fmt.Println(timestampSecret)

		fmt.Printf("Logged into Leetcode as %s", username)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
