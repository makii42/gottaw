package watch

import (
	"github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/output"
	"github.com/makii42/gottaw/pipeline"
	"gopkg.in/urfave/cli.v1"
)

var OneRunCmd = cli.Command{
	Name:    "one",
	Aliases: []string{"1"},
	Usage:   "runs the pipeline once and exits",
	Action:  oneRun,
	Flags:   []cli.Flag{},
}

func oneRun(c *cli.Context) error {
	cfg := config.Load()
	log, err := output.NewLog(cfg)
	if err != nil {
		return err
	}
	pipeline.NewPipeline(nil, log, cfg.Pipeline, func() {
		log.Noticef("Done with run.")
	}).Executor()()
	return nil
}
