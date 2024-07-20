/*
Copyright © 2023 3nthusiast
*/

package cmd

import (
	"backend/log"
	"backend/server"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var l = log.Logger.With().Str("component", "cobra").Logger()

/* FLAGS */
var staticContentPath string
var imageContentPath string
var externalContentPath string
var articleContentPath string
var databasePath string
var srcTranscompilePath string
var dstTranscompilePath string
var tlsCertPath string
var tlsKeyPath string
var isDev *bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "3n-server",
	Short: "3n-server is CLI application, which serves the 3n-app",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		l.Info().Msg("Starting Backend Service")
		server.ServeApp(staticContentPath, imageContentPath, externalContentPath, databasePath, tlsCertPath, tlsKeyPath, isDev)
		l.Info().Msg("Backend Service started successfully")
	},
}

var transcompileCmd = &cobra.Command{
	Use:   "transcompile",
	Short: "Transcompile markdown articles to HTML.",
	Long:  "Iterates over markdown files in the articles repo and transcompile them to the HTML file. It also creates metadata files with title, tags and path properties.",
	Run: func(cmd *cobra.Command, args []string) {
		l.Info().Msg("Transcompiling markdown articles")
		server.TranscompileApp(srcTranscompilePath, dstTranscompilePath)
		l.Info().Msg("Transcompilation completed successfully")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates database with articles from the repo",
	Long: `The initialization of the server is important for the server to load the articles to the database.
	That command should be run in once for a given time period in cron or other task scheduler.`,
	Run: func(cmd *cobra.Command, args []string) {
		l.Info().Msg("Configuration server test initiated")
		server.UpdateApp(articleContentPath, databasePath)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.backend.yaml)")
	rootCmd.Flags().StringVarP(&staticContentPath, "content", "c", "../../3n-app/build", "Provide the path to the static website")
	rootCmd.Flags().StringVarP(&imageContentPath, "img", "", "/app/3n-articles/html/images/", "Provide the path to serve the images for articles")
	rootCmd.PersistentFlags().StringVarP(&articleContentPath, "article", "a", "/app/3n-articles/markdown/", "Provide the path to the articles")
	rootCmd.PersistentFlags().StringVarP(&databasePath, "database", "d", "/app/database.db", "Provide the path to a sqlite database. If the path doesn't exist, it'll create new database.")
	rootCmd.PersistentFlags().StringVarP(&tlsCertPath, "cert", "", "/app/cert.pem", "Provide the path to the TLS certificate")
	rootCmd.PersistentFlags().StringVarP(&tlsKeyPath, "key", "", "/app/key.pem", "Provide the path to the TLS key")
	transcompileCmd.Flags().StringVarP(&srcTranscompilePath, "src", "", "/app/3n-articles/markdown/", "Path to the source markdown articles.")
	transcompileCmd.Flags().StringVarP(&dstTranscompilePath, "dst", "", "/app/3n-articles/html/", "Path to the destination for HTML articles.")

	isDev = rootCmd.PersistentFlags().Bool("dev", false, "Debugging")

	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(transcompileCmd)
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".backend")
	}

	// take into account enviromental variables
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())

		//Viper binds
		viper.BindPFlag("content", rootCmd.Flags().Lookup("content"))
		viper.BindPFlag("img", rootCmd.Flags().Lookup("img"))
		viper.BindPFlag("article", rootCmd.Flags().Lookup("article"))
		viper.BindPFlag("database", rootCmd.Flags().Lookup("database"))
		viper.BindPFlag("cert", rootCmd.Flags().Lookup("cert"))
		viper.BindPFlag("key", rootCmd.Flags().Lookup("key"))
		viper.BindPFlag("dev", rootCmd.Flags().Lookup("dev"))
		viper.BindPFlag("src", transcompileCmd.Flags().Lookup("src"))
		viper.BindPFlag("dst", transcompileCmd.Flags().Lookup("dst"))

		//Viper assings
		staticContentPath = viper.GetString("content")
		imageContentPath = viper.GetString("img")
		articleContentPath = viper.GetString("article")
		databasePath = viper.GetString("database")
		tlsCertPath = viper.GetString("cert")
		tlsKeyPath = viper.GetString("key")
		*isDev = viper.GetBool("dev")
		srcTranscompilePath = viper.GetString("src")
		dstTranscompilePath = viper.GetString("dst")
	}
}
