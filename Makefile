.PHONY: all
all: build

.PHONY: build
build:
	podman build . -t quay.io/jsawaya/terraform-controller-worker

.PHONY: push
push:
	podman push quay.io/jsawaya/terraform-controller-worker

.PHONY: run
run: build
	podman stop terraform-controller
	podman rm terraform-controller
	podman run -d -v "/var/run/docker.sock:/var/run/docker.sock:rw" --name terraform-controller quay.io/jsawaya/terraform-controller-worker

.PHONY: shell
shell: build
	podman stop terraform-controller
	podman rm terraform-controller
	podman run -it -v "/var/run/docker.sock:/var/run/docker.sock:rw" --name terraform-controller quay.io/jsawaya/terraform-controller-worker /bin/sh
