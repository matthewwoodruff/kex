package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "kex",
	Short: "View command examples",
	Args:  cobra.MinimumNArgs(1),
}

const kexFileEnvar = "KEX_FILE"

var output string

func init() {
	kexFileLocation, set := os.LookupEnv(kexFileEnvar)
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "cli", "output mode")

	if set {

		commands, err := parseCommands(kexFileLocation)
		if err != nil {
			panic(fmt.Errorf("error when parsing kex file %v: %w", kexFileLocation, err))
		}

		list := &cobra.Command{
			Use:   "list",
			Short: "List all commands",
			Run: func(cmd *cobra.Command, args []string) {
				listCommands(commands)
			},
		}
		rootCmd.AddCommand(list)

		view := &cobra.Command{
			Use:   "view",
			Short: "View a command",
		}
		rootCmd.AddCommand(view)

		for _, command := range commands {
			view.AddCommand(&cobra.Command{
				Use:   command.Name,
				Short: command.Description,
				Args:  cobra.NoArgs,
				Run:   buildHandlerFunction(command),
			})
		}
	}
}

func listCommands(commands []Command) {
	if output == "md" {
		listCommandsMarkdown(commands)
	} else {
		listCommandsCli(commands)
	}
}

func listCommandsCli(commands []Command) {

}

func listCommandsMarkdown(commands []Command) {
	fmt.Println("# Commands")
	for _, command := range commands {

		if command.Url != "" {
			fmt.Printf("### [%v](%v)\n\n", command.Name, command.Url)
		} else {
			fmt.Printf("### %v\n", command.Name)
		}
		fmt.Printf("%v\n\n", command.Description)
		if command.Notes != "" {
			fmt.Printf("%v\n\n", command.Notes)
		}
		if len(command.Examples) > 0 {
			toMarkdown(command)
		}

		fmt.Println("")
	}
}

func parseCommands(file string) ([]Command, error) {
	bytes, err := ioutil.ReadFile(file)

	if err != nil {
		return []Command{}, fmt.Errorf("failed to read kex file: %w", err)
	}

	var commands []Command

	err = yaml.Unmarshal(bytes, &commands)
	if err != nil {
		return []Command{}, fmt.Errorf("failed to unmarshal kex file: %w", err)
	}

	return commands, nil
}

func buildHandlerFunction(command Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		if output == "md" {
			toMarkdown(command)
		} else {
			printCommandExamples(command)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Command struct {
	Name        string    `yaml:"name"`
	Url         string    `yaml:"url"`
	Description string    `yaml:"description"`
	Notes       string    `yaml:"notes"`
	Examples    []Example `yaml:"examples"`
}

type Example struct {
	Command     string `yaml:"command"`
	Description string `yaml:"description"`
}

func printCommandExamples(command Command) {

	fmt.Printf("%v\n\n", command.Description)
	if command.Notes != "" {
		fmt.Printf("%v\n\n", command.Notes)
	}
	if command.Url != "" {
		fmt.Printf("%v\n\n", command.Url)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetAutoFormatHeaders(false)
	table.SetHeaderLine(false)
	table.SetNoWhiteSpace(true)
	table.SetTablePadding("\t\t")
	table.SetColWidth(150)
	for _, example := range command.Examples {
		table.Append([]string{example.Command, example.Description})
	}
	table.Render()
}

func toMarkdown(command Command) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Example", "Description"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetColWidth(150)

	var data [][]string
	for _, example := range command.Examples {
		data = append(data, []string{fmt.Sprintf("`%v`", example.Command), example.Description})
	}

	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
