/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

/*
write cmd should be simple
water check fanxing
water check fanxing planets flytrap
*/
var checkCmd = &cobra.Command{
	Use:   "check [plant|group] ...",
	Short: "Check a item for today",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var jsonTree, _ = Unmarshal()
		var query = make(map[string]struct{})
		for _, arg := range args {
			query[arg] = struct{}{}
		}
		var items = SelectItems(query, jsonTree)
		WaterItems(items)
		Marshal(jsonTree)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Return dynamic completions based on user's data
		_, validArgs := Unmarshal()
		return validArgs, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
