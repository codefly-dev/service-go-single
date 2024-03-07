package main

import (
	"context"
	"fmt"
	builderv0 "github.com/codefly-dev/core/generated/go/services/builder/v0"
)

type Builder struct {
	*Service
}

func NewBuilder() *Builder {
	return &Builder{
		Service: NewService(),
	}
}
func (s *Builder) Load(ctx context.Context, req *builderv0.LoadRequest) (*builderv0.LoadResponse, error) {
	return s.Builder.LoadError(fmt.Errorf("not implemented"))

}

func (s *Builder) Init(ctx context.Context, req *builderv0.InitRequest) (*builderv0.InitResponse, error) {
	return s.Builder.InitError(fmt.Errorf("not implemented"))

}

func (s *Builder) Update(ctx context.Context, req *builderv0.UpdateRequest) (*builderv0.UpdateResponse, error) {
	return s.Builder.UpdateError(fmt.Errorf("not implemented"))
}

func (s *Builder) Sync(ctx context.Context, req *builderv0.SyncRequest) (*builderv0.SyncResponse, error) {
	return s.Builder.SyncError(fmt.Errorf("not implemented"))

}

func (s *Builder) Build(ctx context.Context, req *builderv0.BuildRequest) (*builderv0.BuildResponse, error) {
	return s.Builder.BuildError(fmt.Errorf("not implemented"))
}

func (s *Builder) Deploy(ctx context.Context, req *builderv0.DeploymentRequest) (*builderv0.DeploymentResponse, error) {
	return s.Builder.DeployError(fmt.Errorf("not implemented"))
}

func (s *Builder) Create(ctx context.Context, req *builderv0.CreateRequest) (*builderv0.CreateResponse, error) {
	return s.Builder.CreateError(fmt.Errorf("not implemented"))
}
