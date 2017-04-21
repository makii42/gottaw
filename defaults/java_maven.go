package defaults

import (
	c "github.com/makii42/gottaw/config"
)

func NewJavaMavenDefault(util *defaultsUtil) *JavaMavenDefault {
	return &JavaMavenDefault{
		util: util,
	}
}

type JavaMavenDefault struct {
	util *defaultsUtil
}

func (j JavaMavenDefault) Name() string {
	return "Java/Maven"
}

func (j JavaMavenDefault) Test(dir string) bool {
	j.util.l.Tracef("testing for %s...\n", j.Name())
	return j.util.fileExists(dir, "pom.xml") &&
		j.util.isExecutable("java") && j.util.isExecutable("javac") &&
		(j.util.isExecutable("mvn") || j.util.isExecutable("mvn.bat"))
}

func (j JavaMavenDefault) Config(dir string) *c.Config {
	return &c.Config{
		Excludes: append(defaultExcludes, "target", "*/target", "*iml"),
		Pipeline: []string{
			"mvn clean test",
		},
	}
}
