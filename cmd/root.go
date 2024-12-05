// Package cmd 命令工具集
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/icdb37/kypd/utils/logx"
)

var rootCmd = &cobra.Command{
	Use:   "kypd",
	Short: "加密/解密",
	Long:  `加密/解密`,
}

// Execute 执行命令
func Execute() {
	if m := os.Getenv("KYPD_INITMASK"); m != "" {
		mask = m
	}
	if l := os.Getenv("KYPD_LOG_LEVEL"); l != "" {
		logx.SetLevel(l)
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// 命令参数
var (
	input    = ""
	output   = ""
	password = ""
	mode     = ""
	mask     = "&?0<GY8B!1w*/f.Zm</9H\"oIbDCO[tc-"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input file")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output file")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", fmt.Sprintf("password length [%d,%d]", minPasswordSize, maxPasswordSize))
	rootCmd.AddCommand(
		cmdEn,
		cmdDe,
	)
}
