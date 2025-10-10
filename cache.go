package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func cacheSave(key string, data []byte) error {
	userCacheLoc, err := os.UserCacheDir()
	if err != nil {
		fmt.Println(redflagStyle.Render(err.Error()))
		return err
	}

	cacheDir := filepath.Join(userCacheLoc, "hhbot")
	if _, err := os.Stat(cacheDir); errors.Is(err, fs.ErrNotExist) {
		fmt.Println(cacheDir, "not exists, try to create...")
		if err := os.Mkdir(cacheDir, 0o700); err != nil {
			fmt.Println(redflagStyle.Render(err.Error()))
			return err
		}
	}

	fname := filepath.Join(cacheDir, key)
	if err := os.WriteFile(fname, data, 0o400); err != nil {
		fmt.Println(redflagStyle.Render(err.Error()))
		return err
	}
	return nil
}

func cacheLoad(key string) ([]byte, error) {
	userCacheLoc, err := os.UserCacheDir()
	if err != nil {
		fmt.Println(redflagStyle.Render(err.Error()))
		return nil, err
	}
	cacheDir := filepath.Join(userCacheLoc, "hhbot")
	fname := filepath.Join(cacheDir, key)
	return os.ReadFile(fname)
}
