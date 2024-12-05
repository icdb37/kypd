package cmd

import (
	"github.com/spf13/cobra"

	"github.com/icdb37/kypd/utils/logx"
)

// cmdDe 加密命令
var cmdEn = &cobra.Command{
	Use:   "en",
	Short: "加密",
	Long:  `加密`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logx.Debugf("input: %s, output: %s, password: %s", input, output, password)
		mode = modeEn
		if err := p.Init(); err != nil {
			return err
		}
		defer p.Close()
		return p.En()
	},
}
