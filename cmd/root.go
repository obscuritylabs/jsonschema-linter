/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/obscuritylabs/jsonschema-linter/lint"
	"github.com/spf13/cobra"
)

var (
	jsonSchema string
	jsonFile   string
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "jsonschema-linter",
	Short:   "Validate JSON files against JSON Schema",
	Long:    `A tool to validate a JSON files using a provided JSON Schema`,
	Version: version + ", commit " + commit + ", built at " + date + " by " + builtBy,
	Args: func(cmd *cobra.Command, args []string) error {
		if jsonSchema == "" {
			return errors.New("you must specify a value for JSON schema")
		}

		if jsonFile == "" {
			return errors.New("you must specify a value for JSON file")
		}
		jsonSchemaPath, err := filepath.Abs(jsonSchema)
		if err != nil {
			return err
		}

		jsonFilePath, err := filepath.Abs(jsonSchema)
		if err != nil {
			return err
		}
		if !isValidPath(jsonSchemaPath) {
			return errors.New("'" + jsonSchemaPath + "' is not a valid path to a json schema file.")
		}

		if !isValidPath(jsonFilePath) {
			return errors.New("'" + jsonFilePath + "' is not a valid path to a json file.")
		}

		jsonSchemaPath = strings.TrimPrefix(jsonSchemaPath, "file://")
		jsonFilePath = strings.TrimPrefix(jsonFilePath, "file://")

		jsonSchema = "file://" + jsonSchemaPath
		jsonFile = "file://" + jsonFilePath

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return lint.Lint(jsonSchema, jsonFile)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&jsonSchema, "schema", "s", "", "JSON Schema (required)")
	rootCmd.MarkFlagRequired("schema")

	rootCmd.Flags().StringVarP(&jsonFile, "json", "j", "", "JSON File (required)")
	rootCmd.MarkFlagRequired("json")

	rootCmd.SetVersionTemplate("")
}

func isValidPath(fp string) bool {
	// https://stackoverflow.com/questions/35231846/golang-check-if-string-is-valid-path
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return true
	}

	return false
}
