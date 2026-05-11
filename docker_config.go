package main

import (
	"bytes"
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

func CheckConfig() []Issue {
	config := exportInfo()
	var output []Issue
	if !config.rootless {
		output = append(output, Issue{
			Severity: SeverityWarning,
			Message: "Docker daemon is not run in rootless mode",
			Fix: "Install docker in rootless mode",
		})
		if !config.userns_remap {
			output = append(output, Issue{
				Severity: SeverityWarning,
				Message: "Userns-remap (User namespace remap) is not enabled while docker is not run in rootless mode.",
				Fix: "Enable userns-remap in your docker config.",
			})
		}
	}
	if !config.docker_content_trust {
		output = append(output, Issue{
			Severity: SeverityInfo,
			Message: "The environment variable DOCKER_CONTENT_TRUST is not enabled, thus docker supports the usage of non-signed images.",
			Fix: "Enable the environment variable DOCKER_CONTENT_TRUST and only use signed images.",
		})
	}
	if !config.cgroupns {
		output = append(output, Issue{
			Severity: SeverityWarning,
			Message: "cgroupns is not enabled in your config",
			Fix: "Enable cgroupns",
		})
	}
	if len(output) != 0 {
		return output
	}
	return nil
}

func exportInfo() Info {
	info := getDockerInfo()
	output := Info{}

	// rootless is recommended
	// cngroups is recommended,
	// userns_remap is recommended if rootless is off
	output.rootless = strings.Contains(info, "rootless")
	output.cgroupns = strings.Contains(info, "cgroupns")
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

	cmd := exec.Command(
		path,
		"context",
		"show",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command(
		path,
		"context",
		"inspect",
		string(bytes.TrimSpace(output[:])),
		"--format",
		`"{{.Endpoints.docker.Host}}"`,
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
