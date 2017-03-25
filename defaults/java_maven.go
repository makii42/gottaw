package defaults

import (
	c "github.com/makii42/gottaw/config"
)

type JavaMavenDefault struct{}

func (j JavaMavenDefault) Name() string {
	return "Java/Maven"
}
func (j JavaMavenDefault) Test(dir string) bool {
	return fileExists(dir, "pom.xml") &&
		isExecutable("java") && isExecutable("javac") &&
		(isExecutable("mvn") || isExecutable("mvn.bat"))
}
func (j JavaMavenDefault) Config() *c.Config {
	return &c.Config{
		Excludes: append(defaultExcludes, "target"),
		Pipeline: []string{
			"mvn clean test",
		},
	}
}
