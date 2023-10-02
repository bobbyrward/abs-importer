package output

import (
	"github.com/drewstinnett/gout/v2"
	"github.com/spf13/cobra"
)

var goutCobraConfig = gout.NewCobraConfig()

func AddOutputFields(cmd *cobra.Command) {
	gout.BindCobra(cmd, goutCobraConfig)
}

func NewGout(cmd *cobra.Command) (*gout.Gout, error) {
	return gout.NewWithCobraCmd(cmd, goutCobraConfig)
}

func Print(cmd *cobra.Command, v interface{}) (err error) {
	gout, err := NewGout(cmd)
	if err != nil {
		return err
	}

	err = gout.Print(v)
	if err != nil {
		return err
	}

	return nil
}

func PrintMulti(cmd *cobra.Command, v ...interface{}) (err error) {
	gout, err := NewGout(cmd)
	if err != nil {
		return err
	}

	err = gout.PrintMulti(v)
	if err != nil {
		return err
	}

	return nil
}
