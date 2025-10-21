package cache

import (
	"encoding/json"
	"os"

	"errors"
	"io/fs"
	"path/filepath"
)

func save(data any, file string, perm os.FileMode) error {
	userCacheLoc, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	hhbotDir := filepath.Join(userCacheLoc, "hhbot")
	if _, err := os.Stat(hhbotDir); errors.Is(err, fs.ErrNotExist) {
		if err := os.Mkdir(hhbotDir, 0o700); err != nil {
			return err
		}
	}

	fname := filepath.Join(hhbotDir, file)
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(f).Encode(data); err != nil {
		return err
	}
	return nil
}

func load(data any, file string) error {
	userCacheLoc, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	hhbotDir := filepath.Join(userCacheLoc, "hhbot")
	if _, err := os.Stat(hhbotDir); errors.Is(err, fs.ErrNotExist) {
		if err := os.Mkdir(hhbotDir, 0o700); err != nil {
			return err
		}
	}
	fname := filepath.Join(hhbotDir, file)
	f, err := os.Open(fname)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}
	return nil
}

func SaveState(state any) error {
	return save(state, "state.json", 0o600)
}

func LoadState(state any) error {
	return load(state, "state.json")
}

func SaveToken(data any) error {
	return save(data, ".accesstoken.json", 0o400)
}

func LoadToken(data any) error {
	return load(data, ".accesstoken.json")
}
