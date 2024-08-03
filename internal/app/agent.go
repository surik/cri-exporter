package app

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	cri "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// Agent defines the interface for an Agent.
type Agent interface {
	GetRuntimeVersion(ctx context.Context) (*cri.VersionResponse, error)
	ListContainers(ctx context.Context) (*cri.ListContainersResponse, error)
	GetStatus(ctx context.Context) (*cri.StatusResponse, error)
	GetRuntimeConfig(ctx context.Context) (*cri.RuntimeConfigResponse, error)
	ListPods(ctx context.Context) (*cri.ListPodSandboxResponse, error)
	ListImages(ctx context.Context) (*cri.ListImagesResponse, error)
	FsInfo(ctx context.Context) (*cri.ImageFsInfoResponse, error)
}

// criAgent represents a CRI agent.
// It is a client to runtime and image services.
type criAgent struct {
	criEndpoint string
	image       cri.ImageServiceClient
	runtime     cri.RuntimeServiceClient
}

// NewAgent creates a new agent.
// It returns an error if the connection to the CRI endpoint fails.
func NewAgent(ctx context.Context, criEndpoint string) (*criAgent, error) {
	conn, err := newConnection(criEndpoint)
	if err != nil {
		return nil, err
	}

	return &criAgent{
		criEndpoint: criEndpoint,
		image:       cri.NewImageServiceClient(conn),
		runtime:     cri.NewRuntimeServiceClient(conn),
	}, nil
}

func newConnection(runtimeEndpoint string) (*grpc.ClientConn, error) {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// TODO: check if the connection is successful

	return grpc.NewClient(runtimeEndpoint, options...)
}

// GetRuntimeVersion returns the version of the container runtime.
func (a *criAgent) GetRuntimeVersion(ctx context.Context) (*cri.VersionResponse, error) {
	return a.runtime.Version(ctx, &cri.VersionRequest{})
}

// ListContainers returns the list of containers in the container runtime.
func (a *criAgent) ListContainers(ctx context.Context) (*cri.ListContainersResponse, error) {
	return a.runtime.ListContainers(ctx, &cri.ListContainersRequest{})
}

// GetStatus returns the status of a container.
func (a *criAgent) GetStatus(ctx context.Context) (*cri.StatusResponse, error) {
	return a.runtime.Status(ctx, &cri.StatusRequest{})
}

// GetRuntimeConfig returns the runtime configuration of the container runtime.
func (a *criAgent) GetRuntimeConfig(ctx context.Context) (*cri.RuntimeConfigResponse, error) {
	return a.runtime.RuntimeConfig(ctx, &cri.RuntimeConfigRequest{})
}

// ListPods returns the list of pods in the container runtime.
func (a *criAgent) ListPods(ctx context.Context) (*cri.ListPodSandboxResponse, error) {
	return a.runtime.ListPodSandbox(ctx, &cri.ListPodSandboxRequest{})
}

// ListImages returns the list of images in the container runtime.
func (a *criAgent) ListImages(ctx context.Context) (*cri.ListImagesResponse, error) {
	return a.image.ListImages(ctx, &cri.ListImagesRequest{})
}

// FsInfo returns the filesystem information of the container runtime.
func (a *criAgent) FsInfo(ctx context.Context) (*cri.ImageFsInfoResponse, error) {
	return a.image.ImageFsInfo(ctx, &cri.ImageFsInfoRequest{})
}
