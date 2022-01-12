package main

// imported packages
import (
    "bufio"
    "fmt"
    "cli-motd/motd/message"
    "os"
    "strings"
    "flag"
)

func main() {
    //Define flags
    var name string
    var greeting string
    var prompt bool
    var preview bool

    // Parse flags
    flag.StringVar(&name, "name", "", "name to use within the message")
    flag.StringVar(&greeting, "greeting", "", "phrase to use within the message")
    flag.BoolVar(&prompt, "prompt", false, "use prompt to input name and greeting")
    flag.BoolVar(&preview, "preview", false, "use preview to output message without writing to /etc/motd")

    flag.Parse()

    // Show usage if flag combination isn`t valid (exit in this scenario)
    if prompt == false && (name == "" || greeting == "") {
        flag.Usage()
        os.Exit(1)
    }

    // Conditionally read from stdin
    if prompt {
        name, greeting = renderPrompt()
    }

    // Optionally print flags and exit if DEBUG variable is set
    if os.Getenv("DEBUG") != "" {
        fmt.Println("Name: ", name)
        fmt.Println("Greeting: ", greeting)
        fmt.Println("Prompt: ", prompt)
        fmt.Println("Preview: ", preview)

        os.Exit(0)
    }

    // Generate message
    message := message.Greeting(name, greeting)

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
