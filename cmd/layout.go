/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ykaseng/hlwm/pkg/exhibiting"
	"github.com/ykaseng/hlwm/pkg/logging"
	"github.com/ykaseng/hlwm/pkg/observing"
)

var set string

// layoutCmd represents the layout command
var layoutCmd = &cobra.Command{
	Use:   "layout",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logging.NewLogger()

		observer := observing.NewService()
		exhibitor := exhibiting.NewService()

		if len(set) > 1 {
			exhibitor.SetLayout(set)
			return
		}

		fmt.Println(observer.Layout())
	},
}

func init() {
	rootCmd.AddCommand(layoutCmd)

	// Here you will define your flags and configuration settings.
	layoutCmd.PersistentFlags().StringVarP(&set, "set", "s", "", "Set hlwm layout")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// layoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// layoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
