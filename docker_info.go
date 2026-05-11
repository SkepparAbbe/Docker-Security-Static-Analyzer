package main

import (
	"os"
	"os/exec"
	"strings"
)

type Info struct {
	rootless             bool
	cgroupns             bool
	userns_remap         bool
	docker_content_trust bool
	docker_socket        string
}

func getDockerInfo() string {
	path, err := exec.LookPath("docker")
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(path, "info")

	output, err := cmd.CombinedOutput()

	if err != nil {
		panic(err)
	}
	return string(output[:])
}

func ExportInfo() Info {
	info := getDockerInfo()
	output := Info{}

	// rootless is recommended
	// cngroups is recommended,
	// userns_remap is recommended if rootless is off
	output.rootless = strings.Contains(info, "rootless")
	output.cgroupns = strings.Contains(info, "cngroups")
	output.userns_remap = strings.Contains(info, "userns")
	output.docker_content_trust = getDockerContentTrust()
	output.docker_socket = getSocketPath()
	return output
}

func getSocketPath() string {
	path, err := exec.LookPath("docker")
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(path, "context show")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command(
		path,
		"context inspect "+string(output[:]),
		" --format \"{{.Endpoints.}}\"",
	)

	output, err = cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	return string(output[:])
}

func getDockerContentTrust() bool {
	_, set := os.LookupEnv("DOCKER_CONTENT_TRUST")
	return set
}
