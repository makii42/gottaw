package pipeline

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/makii42/gottaw/output"
)

type Executor func()

type Pipeline struct {
	pre      Executor
	post     Executor
	commands []string
	wd       string
	log      output.Logger
}

func NewPipeline(preProcess func(), pipeline []string, postProcess func()) *Pipeline {
	return &Pipeline{
		commands: pipeline,
		pre:      preProcess,
		post:     postProcess,
		log:      output.NewLogger(),
	}
}

func (p Pipeline) Executor() Executor {
	return func() {
		start := time.Now()
		if p.pre != nil {
			p.pre()
		}
		for i, commandStr := range p.commands {
			elements := strings.Split(commandStr, " ")
			command, elements := elements[0], elements[1:]
			cmd := exec.Command(command, elements...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if p.wd != "" {
				cmd.Dir = p.wd
			}
			err := cmd.Start()
			if err != nil {
				p.log.Errorf("🚨  (%d@?) ERROR starting '%s': %v", i, commandStr, err)
				return
			}
			pid := cmd.Process.Pid
			p.log.Noticef("♻️  (%d@%d) started '%s'\n", i, pid, commandStr)
			if err := cmd.Wait(); err != nil {
				p.log.Errorf("🚨  (%d@%d) ERROR: %s \n", i, pid, err)
				return
			}
			p.log.Noticef("♻️  (%d@%d) done\n", i, pid)
		}
		dur := time.Since(start)
		p.log.Successf("✅  Pipeline done after %s\n", dur.String())
		if p.post != nil {
			p.post()
		}
	}
}
