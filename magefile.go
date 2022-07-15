//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"planningpoker/infra/dev/mage"
)

const (
	nodeImageTag = "node:14-alpine"
	goImageTag   = "golang:1.15"
)

// Lint runs golang-ci linter.
func Lint() error {
	return sh.RunV("golangci-lint", "run", "-v")
}

// TestUnit runs unit tests with coverage.
func TestUnit() error {
	return sh.RunV("go", "test",
		"-mod=vendor",
		"-count=1",
		"-covermode=atomic",
		"-coverpkg=./...",
		"-coverprofile=coverage.txt",
		"-race",
		"./...")
}

// Test runs all tests with coverage.
func TestAll() error {
	return sh.RunV("go", "test",
		"-mod=vendor",
		"-count=1",
		"-covermode=atomic",
		"-coverpkg=./...",
		"-coverprofile=coverage.txt",
		"-race",
		"-tags=component",
		"./...")
}

// DevFront starts frontend watcher locally with automatic rebuild.
func DevFront() error {
	env, err := mage.GetEnv()
	if err != nil {
		return err
	}

	return sh.RunV("docker", "run",
		"-v", fmt.Sprintf("%s/web:/web", env.WorkingDir),
		"--user", fmt.Sprintf("%s:%s", env.UserID, env.GroupID),
		"--workdir", "/web",
		nodeImageTag, "npm", "run", "build-watch")
}

// CI runs CI pipeline
func CI() {
	mg.SerialDeps(Lint, TestAll)
}

func DockerBuild() error {
	return sh.RunV("docker", "build", "-t", "pp:latest", ".")
}

func DockerRun() error {
	return sh.RunV("docker", "run", "-it", "-p", "8080:8080", "pp:latest")
}
