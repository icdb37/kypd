package cmd

import (
	"github.com/spf13/cobra"

	"github.com/icdb37/kypd/utils/logx"
)

// cmdDe 解密命令
var cmdDe = &cobra.Command{
	Use:   "de",
	Short: "解密",
	Long:  `解密`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mode = modeDe
		logx.Debugf("input: %s, output: %s, password: %s", input, output, password)
		if err := p.Init(); err != nil {
			return err
		}
		defer p.Close()
		return p.De()
	},
}
