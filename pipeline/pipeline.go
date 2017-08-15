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

type (
	command struct {
		Command string
		Shell   string
	}

	pipeline struct {
		commands []string
		wd       string
		log      output.Logger
	}
)

func init() {
	t, err := template.New("script").Parse(scriptTmpl)
	if err != nil {
		panic(err)
	}
	tmpl = t
}

func newPipeline(l output.Logger, commands []string) *pipeline {
	return &pipeline{
		commands: commands,
		log:      l,
	}
}

func (p pipeline) Executor(pre PreFunc, post PostFunc) Executor {
	return func() {
		result := BuildSuccess
		start := time.Now()
		if pre != nil {
			pre()
		}
		defer func() {
			if post != nil {
				post(result)
			}
		}()
		for i, commandStr := range p.commands {

			file, err := ioutil.TempFile("/tmp", "gottaw-")
			if err != nil {
				result = BuildFailure
				panic(err)
			}
			defer os.Remove(file.Name())
			cmdModel := command{Command: commandStr, Shell: "/bin/bash"}
			tmpl.Execute(file, cmdModel)
			if err := file.Close(); err != nil {
				result = BuildFailure
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
				result = BuildFailure
				return
			}
			pid := cmd.Process.Pid
			p.log.Noticef("♻  (%d@%d) started '%s'\n", i, pid, commandStr)
			if err := cmd.Wait(); err != nil {
				p.log.Errorf("🚨  (%d@%d) ERROR: %s \n", i, pid, err)
				result = BuildFailure
				return
			}
			p.log.Noticef("♻  (%d@%d) done\n", i, pid)
		}
		dur := time.Since(start)
		p.log.Successf("✅  Pipeline done after %s\n", dur.String())
	}
}
