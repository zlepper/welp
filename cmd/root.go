/*
 * Copyright © 2018 Rasmus Hansen
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zlepper/welp/internal/app/welp"
	"github.com/zlepper/welp/internal/pkg/models"
	"time"
)

var (
	cfgFile            string
	storageFolderPath  string
	useHttps           bool
	port               int
	tokenDuration      time.Duration
	saveInterval       time.Duration
	databaseFolderPath string
	emailSenderName    string
	emailSenderAddress string
	sendGridApiKey     string
)

const (
	day   = 24 * time.Hour
	month = day * 31
	year  = day * 365
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "welp",
	Short: "A simple server for getting feedback from clients",
	Long: `A very simple server that can help implementing a feedback flow
so clients can easily provide feedback. `,
	Run: func(cmd *cobra.Command, args []string) {
		welp.BindWeb(models.BindWebArgs{
			FolderPath:         storageFolderPath,
			UseHttps:           useHttps,
			Port:               port,
			TokenDuration:      tokenDuration,
			SaveInterval:       saveInterval,
			DatabaseFolderName: databaseFolderPath,
			EmailSenderName:    emailSenderName,
			EmailSenderAddress: emailSenderAddress,
			SendGridApiKey:     sendGridApiKey,
		})
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.welp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	f := rootCmd.Flags()
	f.StringVar(&storageFolderPath, "storageFolderPath", "storage", "Sets the folder welp should storage uploaded files to")

	f.BoolVar(&useHttps, "useHttps", false, "Enable to automatically use https everywhere. Will add http -> https redirect, and enable HSTS. Not compatible with the --port flag")
	f.IntVar(&port, "port", 8080, "Sets the port to host welp on")
	f.DurationVar(&tokenDuration, "tokenDuration", year, "How long a login token should be valid. Shorter times is probably more secure, but longer times makes it easier for users. Defaults to a year.")

	// Flatfile storage options
	f.DurationVar(&saveInterval, "saveInterval", 5*time.Second, "How often the flatFile storage should save changes (such as new feedback, or user changes). Lower values provides better guarantee that data doesn't get lost, but will decrease performance. ")
	f.StringVar(&databaseFolderPath, "databaseFolderPath", "db", "The folder to put database files in.")

	// Email options
	f.StringVar(&emailSenderName, "emailSenderName", "no-reply", "The name that should appear on emails being sent from the system")
	f.StringVar(&emailSenderAddress, "emailSenderAddress", "noreply@noreply.com", "The email address that emails should be sent from. Also used for reply address if people respond to emails.")
	f.StringVar(&sendGridApiKey, "sendGridApiKey", "", "An api key for sendGrid (https://sendgrid.com/). If provided, sendGrid will be used for sending emails.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".welp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".welp")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}
