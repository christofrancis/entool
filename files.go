package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var fileStem = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 ._-]*$`)

func validFileName(input string) (string, error) {
	name := strings.TrimSpace(input)
	if name == "" || strings.ContainsAny(name, `/\:`) {
		return "", errors.New("enter a file name only, not a path")
	}
	if !strings.HasSuffix(strings.ToLower(name), ".txt") {
		if filepath.Ext(name) != "" {
			return "", errors.New("only .txt files are supported")
		}
		name += ".txt"
	}
	stem := name[:len(name)-4]
	if len(name) > 100 || !fileStem.MatchString(stem) || strings.HasSuffix(stem, " ") || strings.HasSuffix(stem, ".") {
		return "", errors.New("use 1-100 characters: English letters, numbers, spaces, ., - or _")
	}
	reserved := strings.ToUpper(strings.SplitN(stem, ".", 2)[0])
	if reserved == "CON" || reserved == "PRN" || reserved == "AUX" || reserved == "NUL" || len(reserved) == 4 && (strings.HasPrefix(reserved, "COM") || strings.HasPrefix(reserved, "LPT")) && reserved[3] >= '1' && reserved[3] <= '9' {
		return "", errors.New("this file name is reserved by Windows")
	}
	return name, nil
}

func executableFolder() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	return filepath.Dir(path), nil
}

func writeEncryptedFile(dir, name, contents string) error {
	normalized, err := validFileName(name)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, normalized)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	_, writeErr := f.WriteString(contents)
	closeErr := f.Close()
	if writeErr != nil || closeErr != nil {
		// Remove partial output if disk fills or a write otherwise fails.
		_ = os.Remove(path)
		return errors.Join(writeErr, closeErr)
	}
	return nil
}

func readEncryptedFile(dir, name string) (string, error) {
	normalized, err := validFileName(name)
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, normalized)
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", errors.New("only regular .txt files inside the EnTool folder are supported")
	}
	if info.Size() > maxCharacters+2 {
		return "", errors.New("encrypted file is too large")
	}
	// Open only after rejecting symlinks. Avoid accepting files outside the folder.
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", errors.New("unable to resolve file location")
	}
	canonicalDir, err := filepath.EvalSymlinks(dir)
	if err != nil || filepath.Dir(resolved) != canonicalDir {
		return "", errors.New("file must be located in the EnTool folder")
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read encrypted file: %w", err)
	}
	return strings.TrimSuffix(strings.TrimSuffix(string(contents), "\n"), "\r"), nil
}

func listTextFiles(dir string) []string {
	matches, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		return nil
	}
	var names []string
	for _, f := range matches {
		info, err := os.Lstat(f)
		if err == nil && info.Mode().IsRegular() {
			names = append(names, filepath.Base(f))
		}
	}
	sort.Strings(names)
	return names
}
