//go:build component
// +build component

package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"planningpoker/infra/dev/mage"
	"planningpoker/internal/infra/repository"
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
	dbHost := fullName + "-postgres"
	_, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:       dbHost,
		Repository: "postgres",
		Tag:        "18.6",
		NetworkID:  env.NetworkID,
		Env: []string{
			"POSTGRES_USER=planningpoker",
			"POSTGRES_PASSWORD=component-test-password",
			"POSTGRES_DB=planningpoker",
		},
	})

	if err != nil {
		return fmt.Errorf("start PostgreSQL: %w", err)
	}

	dsn := fmt.Sprintf("postgres://planningpoker:component-test-password@%s:5432/planningpoker?sslmode=disable", dbHost)
	pool.MaxWait = 30 * time.Second

	if err := pool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		db, err := repository.OpenPostgres(ctx, dsn)
		if err != nil {
			return err
		}
		return db.Close()
	}); err != nil {
		return fmt.Errorf("wait for PostgreSQL: %w", err)
	}

	opts := &dockertest.RunOptions{
		Name:       fullName,
		Repository: pokerContainerName,
		Tag:        env.SessionID,
		NetworkID:  env.NetworkID,
		User:       fmt.Sprintf("%s:%s", env.UserID, env.GroupID),
		Env:        []string{"DATABASE_URL=" + dsn},
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
