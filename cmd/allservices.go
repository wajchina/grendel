package cmd

import (
	"github.com/ubccr/grendel/model"
	"github.com/urfave/cli"
)

func flagExists(flags []cli.Flag, flag cli.Flag) bool {
	for _, f := range flags {
		if flag.GetName() == f.GetName() {
			return true
		}
	}

	return false
}

func NewServeAllCommand() cli.Command {
	flags := make([]cli.Flag, 0)

	for _, cmd := range []cli.Command{NewDHCPCommand(), NewTFTPCommand(), NewAPICommand(), NewDNSCommand()} {
		for _, f := range cmd.Flags {
			if !flagExists(flags, f) {
				flags = append(flags, f)
			}
		}
	}

	return cli.Command{
		Name:        "all",
		Usage:       "Start all services",
		Description: "Start all services",
		Flags:       flags,
		Action:      runAllServices,
	}
}

func runAllServices(c *cli.Context) error {
	staticBooter, err := model.NewStaticBooter(c.String("static-hosts"), c.String("kernel"), c.StringSlice("initrd"), c.String("cmdline"), c.String("liveimg"))
	if err != nil {
		return err
	}

	DB = staticBooter

	errs := make(chan error, 4)

	go func() { errs <- runDHCP(c) }()
	go func() { errs <- runTFTP(c) }()
	go func() { errs <- runAPI(c) }()
	go func() { errs <- runDNS(c) }()

	err = <-errs
	return err
}
