/*
Copyright © 2022 Oleksandr Patuk

*/
package cmd

import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var cfgFile string
var name string
var greeting string
var preview bool
var prompt bool
var debug bool = false

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "cli-motd-cobra",
    Short: "A utility to customize the Message of the Day (MOTD)",
    Long: ``,
    // Uncomment the following line if your bare application
    // has an action associated with it:
    Run: func(cmd *cobra.Command, args []string) {
        // If no arguments passed, show usage
        if !prompt && (name == "" || greeting == "") {
            cmd.Usage()
            os.Exit(0)
        }

        // Optionally print flags and exit if DEBUG variable is set
        if debug {
            fmt.Println("Name: ", name)
            fmt.Println("Greeting: ", greeting)
            fmt.Println("Prompt: ", prompt)
            fmt.Println("Preview: ", preview)
            os.Exit(0)
        }

        // Conditionally read from stdin
        if prompt {
            name, greeting = renderPrompt()
        }

        //Generate message
        message := buildMessage(name, greeting)

        // Either preview message or write to file
        if preview {
            fmt.Println(message)
        } else {
            // Write content
            f, err := os.OpenFile("/etc/motd", os.O_WRONLY, 0644)

            if err != nil {
                fmt.Println("Error: Unable to open /etc/motd")
                os.Exit(1)
            }

            defer f.Close()

            _, err = f.Write([]byte(message))

            if err != nil {
                fmt.Println("Error: Unable to write to /etc/motd")
                os.Exit(1)
            }
        }
    },
}

func buildMessage (name, message string) (salutation string) {
    salutation = fmt.Sprintf("%s, %s", message, name)
    return
}

func renderPrompt() (name, greeting string) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Your Greeting: ")
    greeting, _ = reader.ReadString('\n')
    greeting = strings.TrimSpace(greeting)

    fmt.Print("Your Name: ")
    name, _ = reader.ReadString('\n')
    name = strings.TrimSpace(name)
    return
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    // Here you will define your flags and configuration settings.
    // Cobra supports persistent flags, which, if defined here,
    // will be global for your application.

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.cli-motd-cobra.yaml)")

    // Cobra also supports local flags, which will only run
    // when this action is called directly.
    rootCmd.Flags().StringVarP(&name, "name", "n", "", "Name to use in the message")
    rootCmd.Flags().StringVarP(&greeting, "greeting", "g", "", "Greeting to use in the message")
    rootCmd.Flags().BoolVarP(&preview, "preview", "v", false, "Preview message instead of writing to /etc/motd")
    rootCmd.Flags().BoolVarP(&prompt, "prompt", "p", false, "Prompt for name and greeting")

    if os.Getenv("DEBUG") != "" {
        debug = true
    }
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    if cfgFile != "" {
        // Use config file from the flag.
        viper.SetConfigFile(cfgFile)
    } else {
        // Search config in the current directory with name ".cli-motd-cobra" (without extension).
        viper.AddConfigPath("./")
        viper.SetConfigName(".cli-motd-cobra")
    }

    viper.AutomaticEnv() // read in environment variables that match

    // If a config file is found, read it in.
    if err := viper.ReadInConfig(); err == nil {
        fmt.Println("Using config file:", viper.ConfigFileUsed())
    }
}
