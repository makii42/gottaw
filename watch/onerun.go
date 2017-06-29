package watch

import (
	"github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/output"
	"gopkg.in/urfave/cli.v1"
)

var OneRunCmd = cli.Command{
	Name:    "one",
	Aliases: []string{"1"},
	Usage:   "runs the pipeline once and exits",
	Action:  OneRun,
	Flags:   []cli.Flag{},
}

func OneRun(c *cli.Context) {
	watchCfg = config.Setup(c.GlobalString("config"))
	log = output.NewLogger(output.TRACE, watchCfg)
	executePipeline(nil, watchCfg.Pipeline, func() {
		log.Noticef("Done with run.")
	})()

}
