/*
Copyright © 2021 Ka Seng <me@ykaseng.com>

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/ykaseng/hlwm/pkg/exhibiting"
	"github.com/ykaseng/hlwm/pkg/logging"
	"github.com/ykaseng/hlwm/pkg/observing"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status observes hlwm tag state changes",
	Long: `Status observes hlwm tag state changes, outputs changes to stdout 
	and signals eww to exhibit changes. Errors are piped through to notify-send
	using logrus.`,
	Run: func(cmd *cobra.Command, args []string) {
		logging.NewLogger()

		exhibitor := exhibiting.NewService()
		observer := observing.NewService()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		tcs := observer.TagChangeEvent(ctx)
		fs := exhibitor.FlashWidget(ctx, "statusbar", time.Second)
		for {
			defer close(fs)
			select {
			case <-ctx.Done():
				return
			case <-tcs:
				fmt.Println(observer.TagStatus())
				fs <- struct{}{}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
