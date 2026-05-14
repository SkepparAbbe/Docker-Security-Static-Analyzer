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

// Finds and checks the user's docker config. Searches for keywwords and reports
// issues for individual or combinations of unsafe configurations. 
func CheckConfig() []Issue {
	config := exportInfo()
	var output []Issue
	if !config.rootless {
		output = append(output, Issue{
			Severity: SeverityWarning,
			Problem: "Docker daemon is not run in rootless mode",
			Description: "The Daemon has complete root access and any adversary inside the container who gains access over it gains full control.",
			Fix: "Install docker in rootless mode",
		})
		if !config.userns_remap {
			output = append(output, Issue{
				Severity: SeverityWarning,
				Problem: "Userns-remap (User namespace remap) is not enabled while docker is not run in rootless mode.",
				Description: "Userns-remap maps the root of the container to a regular user on the host, giving an adversary who gains root access in the container less privileges.",
				Fix: "Enable userns-remap in your docker config.",
			})
		}
	}
	if !config.docker_content_trust {
		output = append(output, Issue{
			Severity: SeverityInfo,
			Problem: "The environment variable DOCKER_CONTENT_TRUST is not enabled",
			Description: "This means that you are not verifying the integrity and publisher of the images you pull. This can lead to supply chain attacks.",
			Fix: "Enable the environment variable DOCKER_CONTENT_TRUST and only use signed images.",
		})
	}
	if !config.cgroupns {
		output = append(output, Issue{
			Severity: SeverityWarning,
			Problem: "cgroupns is not enabled in your config",
			Description: "A user inside the container cannot see the host's global cgroup tree if cgroupns is enabled.",
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
