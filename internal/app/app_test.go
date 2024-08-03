package app_test

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/surik/cri-exporter/internal/app"
	cri "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func TestDescribeByCollect(t *testing.T) {
	metrics := app.NewMetrics(&mockAgent{}, "test")
	assert.NotNil(t, metrics)

	reg := prometheus.NewRegistry()
	err := reg.Register(metrics)
	assert.NoError(t, err)

	_, err = reg.Gather()
	assert.NoError(t, err)
}

type mockAgent struct{}

func (a *mockAgent) GetRuntimeVersion(ctx context.Context) (*cri.VersionResponse, error) {
	return &cri.VersionResponse{}, nil
}

func (a *mockAgent) ListContainers(ctx context.Context) (*cri.ListContainersResponse, error) {
	return &cri.ListContainersResponse{}, nil
}

func (a *mockAgent) GetStatus(ctx context.Context) (*cri.StatusResponse, error) {
	return &cri.StatusResponse{Status: &cri.RuntimeStatus{Conditions: []*cri.RuntimeCondition{
		{Type: "Ready", Status: true},
	}}}, nil
}

func (a *mockAgent) GetRuntimeConfig(ctx context.Context) (*cri.RuntimeConfigResponse, error) {
	return &cri.RuntimeConfigResponse{}, nil
}

func (a *mockAgent) ListPods(ctx context.Context) (*cri.ListPodSandboxResponse, error) {
	return &cri.ListPodSandboxResponse{}, nil
}

func (a *mockAgent) ListImages(ctx context.Context) (*cri.ListImagesResponse, error) {
	return &cri.ListImagesResponse{}, nil
}

func (a *mockAgent) FsInfo(ctx context.Context) (*cri.ImageFsInfoResponse, error) {
	return &cri.ImageFsInfoResponse{}, nil
}
