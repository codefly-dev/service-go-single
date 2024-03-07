package main

import (
	"context"
	"github.com/codefly-dev/core/agents/services"
	"github.com/codefly-dev/core/configurations"
	basev0 "github.com/codefly-dev/core/generated/go/base/v0"
	"github.com/codefly-dev/core/wool"

	agentv0 "github.com/codefly-dev/core/generated/go/services/agent/v0"

	"github.com/codefly-dev/core/agents/helpers/code"
	golanghelpers "github.com/codefly-dev/core/agents/helpers/go"
	runtimev0 "github.com/codefly-dev/core/generated/go/services/runtime/v0"
)

type Runtime struct {
	*Service

	// internal
	runner *golanghelpers.Runner

	Environment          *basev0.Environment
	EnvironmentVariables *configurations.EnvironmentVariableManager

	NetworkMappings []*basev0.NetworkMapping
	Port            int
}

func NewRuntime() *Runtime {
	return &Runtime{
		Service: NewService(),
	}
}

func (s *Runtime) Load(ctx context.Context, req *runtimev0.LoadRequest) (*runtimev0.LoadResponse, error) {
	defer s.Wool.Catch()
	ctx = s.Wool.Inject(ctx)

	s.Configuration = &configurations.Service{}

	s.Settings.Watch = true

	err := s.Base.HeadlessLoad(ctx, req.Identity, s.Settings)
	if err != nil {
		return s.Base.Runtime.LoadError(err)
	}

	s.Wool.Focus("location", wool.Field("location", s.Location))

	requirements.Localize(s.Location)

	s.Environment = req.Environment

	s.EnvironmentVariables = s.LoadEnvironmentVariables(req.Environment)

	if s.Settings.Watch {
		s.Wool.Debug("setting up code watcher", wool.Field("configurations", requirements))
		conf := services.NewWatchConfiguration(requirements)
		err = s.SetupWatcher(ctx, conf, s.EventHandler)
		if err != nil {
			s.Wool.Warn("error in watcher", wool.ErrField(err))
		}
	}

	s.EnvironmentVariables = configurations.NewEnvironmentVariableManager()

	return s.Base.Runtime.LoadResponse()
}

func (s *Runtime) Init(ctx context.Context, req *runtimev0.InitRequest) (*runtimev0.InitResponse, error) {
	defer s.Wool.Catch()
	ctx = s.Wool.Inject(ctx)

	s.NetworkMappings = req.ProposedNetworkMappings

	runner, err := golanghelpers.NewRunner(ctx, s.Location)
	if err != nil {
		return s.Runtime.InitError(err)
	}

	// Stop before replacing the runner
	if s.runner != nil {
		err = s.runner.Stop()
		if err != nil {
			return s.Runtime.InitError(err)
		}
	}
	s.runner = runner

	s.runner.WithDebug(s.Settings.Debug)
	s.runner.WithRaceConditionDetection(s.Settings.WithRaceConditionDetectionRun)
	s.runner.WithRequirements(requirements)

	// Output to wool
	s.runner.WithOut(s.Wool)

	s.Wool.Debug("runner init started")
	err = s.runner.Init(ctx)
	if err != nil {
		s.Wool.Error("cannot init the go runner", wool.ErrField(err))
		return s.Runtime.InitError(err)
	}
	s.Wool.Debug("runner init done")
	s.Ready()

	s.Wool.Info("successful init of runner")

	return s.Runtime.InitResponse(s.NetworkMappings)
}

func (s *Runtime) Start(ctx context.Context, req *runtimev0.StartRequest) (*runtimev0.StartResponse, error) {
	defer s.Wool.Catch()
	ctx = s.Wool.Inject(ctx)

	runningContext := s.Wool.Inject(context.Background())

	s.Wool.Debug("starting runner")
	err := s.runner.Start(runningContext)
	if err != nil {
		return s.Runtime.StartError(err, wool.Field("in", "runner"))
	}
	s.Wool.Debug("runner started successfully")

	return s.Runtime.StartResponse()
}

func (s *Runtime) Information(ctx context.Context, req *runtimev0.InformationRequest) (*runtimev0.InformationResponse, error) {
	return s.Base.Runtime.InformationResponse(ctx, req)
}

func (s *Runtime) Stop(ctx context.Context, req *runtimev0.StopRequest) (*runtimev0.StopResponse, error) {
	defer s.Wool.Catch()

	s.Wool.Debug("stopping service")
	err := s.runner.Stop()
	if err != nil {
		return s.Runtime.StopError(err)
	}

	err = s.Base.Stop()
	if err != nil {
		return s.Runtime.StopError(err)
	}
	return s.Runtime.StopResponse()
}

func (s *Runtime) Communicate(ctx context.Context, req *agentv0.Engage) (*agentv0.InformationRequest, error) {
	return s.Base.Communicate(ctx, req)
}

/* Details

 */

func (s *Runtime) EventHandler(event code.Change) error {
	s.Wool.Debug("detected change requiring re-build", wool.Field("path", event.Path))
	s.Runtime.DesiredInit()
	return nil
}
