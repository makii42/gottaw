package pipeline

import (
	"sync"

	"github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/output"
)

type (
	// BuildResult describes the outcome of a build, either success and failure
	BuildResult int

	// Executor runs the actual pipeline commands
	Executor func()
	// PreFunc will be executed before before a pipeline kicks in
	PreFunc func()
	// PostFunc will be executed after a pipeline finished, and receive
	// the build result as parameter
	PostFunc func(BuildResult)

	// Builder offers a thing that can create an executor to build a pipeline
	Builder interface {
		Executor(PreFunc, PostFunc) (Executor, error)
	}
	builder struct {
		cfg *config.Config
		log output.Logger
		mux sync.Mutex
	}
)

const (
	// BuildSuccess indicates a successful build
	BuildSuccess BuildResult = iota
	// BuildFailure indicates a failed build
	BuildFailure
)

// NewBuilder creates a new builder.
func NewBuilder(cfg *config.Config, log output.Logger) Builder {
	return &builder{
		cfg: cfg,
		log: log,
	}
}

// Executor will create an executor for the pipeline
func (b *builder) Executor(pre PreFunc, post PostFunc) (Executor, error) {
	pipeline := newPipeline(b.log, b.cfg.Pipeline)
	preWrap := func() {
		b.mux.Lock()
		if pre != nil {
			pre()
		}
	}
	postWrap := func(r BuildResult) {
		b.mux.Unlock()
		if post != nil {
			post(r)
		}
	}
	return pipeline.Executor(preWrap, postWrap), nil
}
