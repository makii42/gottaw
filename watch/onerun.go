package watch

import (
	"github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/output"
	"github.com/makii42/gottaw/pipeline"
	"gopkg.in/urfave/cli.v1"
)

// OneRunCmd is the command to run a pipeline once.
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
	builder := pipeline.NewBuilder(cfg, log)
	executor, err := builder.Executor(nil, func(r pipeline.BuildResult) {
		var resMsg string
		if r == pipeline.BuildSuccess {
			resMsg = "succeeded"
		} else {
			resMsg = "failed"
		}
		log.Noticef("The build %s.\n", resMsg)
	})
	if err != nil {
		return err
	}
	executor()
	return nil
}
