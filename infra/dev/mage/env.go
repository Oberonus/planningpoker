package mage

import (
	"errors"
	"os"

	"github.com/magefile/mage/sh"
)

// Env is a storage of testing environment variables.
type Env struct {
	UserID     string
	GroupID    string
	WorkingDir string
	NetworkID  string
	SessionID  string
}

var env *Env

// GetEnv returns current testing environment variables.
func GetEnv() (*Env, error) {
	// return cached value if possible
	if env != nil {
		return env, nil
	}

	// user.Current() does not work in container, let's go hardcore.
	uid, err := sh.Output("id", "-u")
	if err != nil {
		return nil, err
	}
	gid, err := sh.Output("id", "-g")
	if err != nil {
		return nil, err
	}

	workDir := os.Getenv("REPO_DIRECTORY")
	if workDir == "" {
		if workDir, err = os.Getwd(); err != nil {
			return nil, err
		}
	}

	sid := os.Getenv("MAGE_SESSION_ID")
	if sid == "" {
		return nil, errors.New("MAGE_SESSION_ID is not set")
	}

	nid := os.Getenv("MAGE_NETWORK_ID")
	if sid == "" {
		return nil, errors.New("MAGE_NETWORK_ID is not set")
	}

	env = &Env{
		UserID:     uid,
		GroupID:    gid,
		WorkingDir: workDir,
		SessionID:  sid,
		NetworkID:  nid,
	}

	return env, nil
}
