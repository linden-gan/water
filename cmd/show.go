/*
Copyright © 2025 ganlinden@gmail.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

/*
By default, in the same order as the json file
water show xiaoshu
water show fanxing
water show fanxing xiaoshu flytrap

// Sort based on frequency
water show -f
water show -f plants

// Sort based on urgency
water show -u
water show -u plants
*/
var showCmd = &cobra.Command{
	Use:   "show [plant|group] ...",
	Short: "Show watering schedule of plants or groups",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var jsonTree, _ = Unmarshal()
		var query = make(map[string]struct{})
		for _, arg := range args {
			query[arg] = struct{}{}
		}
		if len(query) == 0 {
			query["all"] = struct{}{}
		}
		var items = SelectItems(query, jsonTree)
		var cs, earliestDate = Json2CliStructs(items)
		Print(cs, earliestDate)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Return dynamic completions based on user's data
		_, validArgs := Unmarshal()
		return validArgs, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
