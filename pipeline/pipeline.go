package pipeline

import (
	"os"
	"os/exec"
	"time"

	"io/ioutil"
	"text/template"

	"github.com/makii42/gottaw/output"
)

const scriptTmpl = `#!{{ .Shell }}
{{ .Command }}`

var (
	tmpl *template.Template
)

type Executor func()

type command struct {
	Command string
	Shell   string
}

type Pipeline struct {
	pre      Executor
	post     Executor
	commands []string
	wd       string
	log      output.Logger
}

func init() {
	t, err := template.New("script").Parse(scriptTmpl)
	if err != nil {
		panic(err)
	}
	tmpl = t
}

func NewPipeline(preProcess func(), l output.Logger, pipeline []string, postProcess func()) *Pipeline {
	return &Pipeline{
		commands: pipeline,
		pre:      preProcess,
		post:     postProcess,
		log:      l,
	}
}

func (p Pipeline) Executor() Executor {
	return func() {
		start := time.Now()
		if p.pre != nil {
			p.pre()
		}
		for i, commandStr := range p.commands {

			file, err := ioutil.TempFile("/tmp", "gottaw-")
			if err != nil {
				panic(err)
			}
			defer os.Remove(file.Name())
			cmdModel := command{Command: commandStr, Shell: "/bin/bash"}
			tmpl.Execute(file, cmdModel)
			if err := file.Close(); err != nil {
				panic(err)
			}

			cmd := exec.Command("/bin/bash", file.Name())
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if p.wd != "" {
				cmd.Dir = p.wd
			}
			if err := cmd.Start(); err != nil {
				p.log.Errorf("🚨  (%d@?) ERROR starting '%s': %v", i, commandStr, err)
				return
			}
			pid := cmd.Process.Pid
			p.log.Noticef("♻  (%d@%d) started '%s'\n", i, pid, commandStr)
			if err := cmd.Wait(); err != nil {
				p.log.Errorf("🚨  (%d@%d) ERROR: %s \n", i, pid, err)
				return
			}
			p.log.Noticef("♻  (%d@%d) done\n", i, pid)
		}
		dur := time.Since(start)
		p.log.Successf("✅  Pipeline done after %s\n", dur.String())
		if p.post != nil {
			p.post()
		}
	}
}
