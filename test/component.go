//go:build component
// +build component

package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"planningpoker/infra/dev/mage"
)

const (
	pokerContainerName = "planning-poker"
)

var (
	pokerHost string
	pokerPort string
)

// StartPokerContainer starts a planning poker container.
func StartPokerContainer() error {
	env, err := mage.GetEnv()
	if err != nil {
		return err
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		return err
	}

	imageName := pokerContainerName + ":" + env.SessionID
	err = pool.Client.BuildImage(docker.BuildImageOptions{
		Name:         imageName,
		Dockerfile:   "Dockerfile",
		OutputStream: ioutil.Discard,
		ContextDir:   "..",
	})
	if err != nil {
		return fmt.Errorf("build image %s: %w", imageName, err)
	}

	fullName := env.SessionID + "-" + pokerContainerName
	opts := &dockertest.RunOptions{
		Name:       fullName,
		Repository: pokerContainerName,
		Tag:        env.SessionID,
		NetworkID:  env.NetworkID,
		User:       fmt.Sprintf("%s:%s", env.UserID, env.GroupID),
	}

	_, err = pool.RunWithOptions(opts)
	if err != nil {
		return err
	}

	pokerHost = fullName
	pokerPort = "8080"

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = time.Second * 20
	if err = pool.Retry(func() error {
		resp, err := http.Get(fmt.Sprintf("http://%s:%s/alive", pokerHost, pokerPort))
		if err != nil {
			return err
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("got status code: %d", resp.StatusCode)
		}

		return nil
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	return nil
}

func isRunningInDockerContainer() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}
