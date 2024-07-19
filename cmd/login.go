package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/juju/persistent-cookiejar"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
// TODO: scanf login details and make password invisible (and 2fa)
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
		password := args[1]

		jar, _ := cookiejar.New(nil)

		client := http.Client{
			Jar: jar,
		}

		// leetcodeLogin := "https://leetcode.com/accounts/github/login/?next=%2F"
		githubLogin := "https://github.com/login"
		githubSession := "https://github.com/session"

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
		authToken := authRegex.FindStringSubmatch(bodyString)

		// idRegex, _ := regexp.Compile(`name="ga_id" value="(.*?)"`)
		// gaId := idRegex.FindStringSubmatch(bodyString)
		// if len(gaId) == 0 {
		// 	gaId = []string{"", ""}
		// }

		timestampRegex, _ := regexp.Compile(`name="timestamp" value="(.*?)"`)
		timestamp := timestampRegex.FindStringSubmatch(bodyString)

		timestampSecretRegex, _ := regexp.Compile(`name="timestamp_secret" value="(.*?)"`)
		timestampSecret := timestampSecretRegex.FindStringSubmatch(bodyString)

		if len(authToken) == 0 || len(timestamp) == 0 || len(timestampSecret) == 0 {
			fmt.Println("Failed to find required values")
			return
		}

		signinURI := "?x=Sign%20in"
		form := url.Values{}
		form.Add("login", username)
		form.Add("password", password)
		form.Add("authenticity_token", authToken[1])
		form.Add("commit", signinURI)
		form.Add("ga_id", "")
		form.Add("webauthn-support", "supported")
		form.Add("webauthn-iuvpaa-support", "unsupported")
		form.Add("return_to", "")
		form.Add("required_field", "")
		form.Add("timestamp", timestamp[1])
		form.Add("timestamp_secret", timestampSecret[1])

		req, err := http.NewRequest("POST", githubSession, strings.NewReader(form.Encode()))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		sessionRes, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		defer sessionRes.Body.Close()

		fmt.Println(sessionRes.StatusCode)
		sbodyBytes, err := io.ReadAll(sessionRes.Body)
		if err != nil {
			print(err)
		}
		sbodyString := string(sbodyBytes)
		fmt.Println(sbodyString)
		// Save cookie jar afterwards https://akimon658.github.io/en/p/2022/go-local-cookie/

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
