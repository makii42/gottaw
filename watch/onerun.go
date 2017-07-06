package watch

import (
	"github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/output"
	"gopkg.in/urfave/cli.v1"
	pipeline2 "github.com/makii42/gottaw/pipeline"
)

var OneRunCmd = cli.Command{
	Name:    "one",
	Aliases: []string{"1"},
	Usage:   "runs the pipeline once and exits",
	Action:  oneRunFactory,
	Flags:   []cli.Flag{},
}

func oneRunFactory(cfg *config.Config, out *output.Output) cli.ActionFunc {

	return func(c *cli.Context) error {
		log := output.NewLogger()

		pipeline := pipeline2.NewPipeline(nil, cfg.Pipeline, func() {
			log.Noticef("Done with run.")
		})
		pipeline.Executor()()
		return nil
	}
}
