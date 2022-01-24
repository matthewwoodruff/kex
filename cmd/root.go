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

func init() {
	kexFileLocation, set := os.LookupEnv(kexFileEnvar)

	if set {
		commands, err := parseCommands(kexFileLocation)
		if err != nil {
			panic(fmt.Errorf("error when parsing kex file %v: %w", kexFileLocation, err))
		}

		for _, command := range commands {
			rootCmd.AddCommand(&cobra.Command{
				Use:   command.Name,
				Short: command.Description,
				Args:  cobra.NoArgs,
				Run:   buildHandlerFunction(command),
			})
		}
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
		printCommandExamples(command)
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
