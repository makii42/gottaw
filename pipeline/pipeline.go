package pipeline

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

type Executor func()

type Pipeline struct {
	pre      Executor
	post     Executor
	commands []string
	wd       string
}

func NewPipeline(preProcess func(), wd string, pipeline []string, postProcess func()) *Pipeline {

	return &Pipeline{
		commands: pipeline,
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
				log.Errorf("🚨  (%d@?) ERROR starting '%s': %v", i, commandStr, err)
				return
			}
			pid := cmd.Process.Pid
			log.Noticef("♻️  (%d@%d) started '%s'\n", i, pid, commandStr)
			if err := cmd.Wait(); err != nil {
				log.Errorf("🚨  (%d@%d) ERROR: %s \n", i, pid, err)
				return
			}
			log.Noticef("♻️  (%d@%d) done\n", i, pid)
		}
		dur := time.Since(start)
		log.Successf("✅  Pipeline done after %s\n", dur.String())
		if p.post != nil {
			p.post()
		}
	}
}
