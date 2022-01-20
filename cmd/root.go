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
}

func init() {
	bytes, err := ioutil.ReadFile("commands.yaml")

	if err != nil {
		panic(err)
	}

	var commands []Command

	err = yaml.Unmarshal(bytes, &commands)
	if err != nil {
		panic(err)
	}

	for _, command := range commands {
		rootCmd.AddCommand(&cobra.Command{
			Use:   command.Name,
			Short: command.Description,
			Args:  cobra.NoArgs,
			Run:   sad(command),
		})
	}
}

func sad(command Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		root(command)
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
	Examples    []Example `yaml:"examples"`
}

type Example struct {
	Command     string `yaml:"command"`
	Description string `yaml:"description"`
}

func root(command Command) {

	fmt.Printf("%v\n\n", command.Description)
	fmt.Printf("%v\n\n", command.Url)
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
