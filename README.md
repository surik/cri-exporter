# CRI-Exporter

This is a simple tool to export data from the CRI API as Prometheus metrics.

You can test it on MacOS with Docker Desktop. You need to have Kubernetes enabled and running. 
Then you can run the following commands:

    $ make docker run-docker

Then you can access the metrics on http://localhost:8080/metrics

**NOTE:** For local tests, this is necessary to have Kubernetes enabled in Docker Desktop as this will allow to mount `cri-dockerd.sock` into the container.

You may try to test it with `CRI-O` or `containerd`. PRs are welcome.