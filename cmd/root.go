package cmd

import (
	"errors"
	"fmt"
	"memlogy/database"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "memlogy",
	Short: "memlogy allows you to easily track your daily history",
	Long:  `your notes, what you did, etc. Write anything you want easily and easily track your history.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { fmt.Println("Hello CLI") },
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
	//fmt.Println("flag")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.memlogy.yaml)")

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
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".memlogy" (without extension).
		viper.AddConfigPath(home + "/")
		viper.SetConfigName(".memlogy")
		viper.SetConfigType("yaml")

		fileName := home + "/.memlogy"
		if _, err := os.Stat(home + "/.memlogy"); errors.Is(err, os.ErrNotExist) {
			fmt.Println("File not found")
			file, err := os.Create(fileName)
			file.Close()

			if err != nil {
				fmt.Println(err)
			}
			viper.SetDefault("db-initialized", true)
			viper.WriteConfigAs(fileName)
			database.Migrate()
		}
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Config File not found")
		} else {
			// Config file was found but another error was produced
			fmt.Println("err" + err.Error())
		}
	}

	//fmt.Println("read config :" + viper.GetString("db-initialized"))
}
