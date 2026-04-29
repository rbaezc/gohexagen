package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/rbaezc/gohexagen/commands"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gohexagen",
		Short: "Hexagonal Architecture Scaffolder for Go",
	}

	var initCmd = &cobra.Command{
		Use:   "init [name]",
		Short: "Initializes a new Hexagonal Go project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			fmt.Printf("🛡️ Initializing Go project: %s\n\n", projectName)

			framework := ""
			promptFw := &survey.Select{
				Message: "Select your Application Framework (Primary Adapter):",
				Options: []string{
					"Fiber (Recommended - Ultra Fast)", 
					"Gin (Classic & Battle-tested)", 
					"Echo (Elegant & Scalable)", 
					"None (I will wire my own handlers)",
				},
			}
			if err := survey.AskOne(promptFw, &framework); err != nil {
				fmt.Println("Error reading input:", err)
				return
			}

			orm := ""
			promptOrm := &survey.Select{
				Message: "Select your Database Tool (Secondary Adapter):",
				Options: []string{
					"GORM (Default - Feature Rich)", 
					"Ent (Modern Graph-based ORM)", 
					"SQLx (Lightweight Raw SQL)", 
					"Custom (I will map the interfaces myself)",
				},
			}
			if err := survey.AskOne(promptOrm, &orm); err != nil {
				fmt.Println("Error reading input:", err)
				return
			}

			fmt.Printf("\n🚀 Selected Framework: %s\n", framework)
			fmt.Printf("🚀 Selected ORM: %s\n", orm)
			
			if err := commands.InitProject(projectName, framework, orm); err != nil {
				fmt.Println("❌ Error scaffolding project:", err)
				return
			}
			fmt.Println("\n✅ Project successfully initialized!")
		},
	}

	var genCmd = &cobra.Command{
		Use:   "gen-resource [entity] [fields...]",
		Short: "Generates a new vertical slice",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			entityName := args[0]
			fields := args[1:]
			if err := commands.GenerateResource(entityName, fields); err != nil {
				fmt.Println("❌ Error generating resource:", err)
				return
			}
			fmt.Println("\n✅ Resource successfully generated!")
		},
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
