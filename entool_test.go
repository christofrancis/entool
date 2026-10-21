package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Created with the original Python EnTool format: scrypt N=32768,r=8,p=1,
// 16-byte salt, 12-byte nonce, AES-256-GCM AAD="ET1", and base64url filler.
const pythonET1Fixture = "ET1AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGx_-C0mr1n4C-hl4faZnEa1wIRRbxWuP0FFdgmBmNzDawepKMFOxsdzHY1noG8aQtjgpvo7VlvdIwl0c8w==ab"

func TestLegacyPythonFile(t *testing.T) {
	got, err := decryptSecret(pythonET1Fixture, "pass phrase ★")
	if err != nil || got != "Tajne hasło 🗝️" {
		t.Fatalf("Python ET1 compatibility failed: %q, %v", got, err)
	}
	if _, err := decryptSecret(pythonET1Fixture, "wrong"); !errors.Is(err, errInvalidEncryptedFile) {
		t.Fatalf("wrong key should not decrypt: %v", err)
	}
}

func TestExactLengthRoundTrip(t *testing.T) {
	for _, text := range []string{"abc", "żółw 🗝️", "1234 中文"} {
		for extra := 0; extra < 8; extra++ {
			n := minimumCharacters(text) + extra
			encrypted, err := encryptSecret(text, "a very long passphrase", n)
			if err != nil {
				t.Fatalf("encrypt len %d: %v", n, err)
			}
			if len(encrypted) != n {
				t.Fatalf("requested %d characters but got %d", n, len(encrypted))
			}
			got, err := decryptSecret(encrypted, "a very long passphrase")
			if err != nil || got != text {
				t.Fatalf("decrypt len %d: %q %v", n, got, err)
			}
			if len(encrypted) < 4 || !strings.HasPrefix(encrypted, formatVersion) {
				t.Fatal("wrong format prefix")
			}
		}
	}
}

func TestInvalidLengthsAndTampering(t *testing.T) {
	if _, err := encryptSecret("abc", "password", minimumCharacters("abc")-1); err == nil {
		t.Fatal("too small accepted")
	}
	if _, err := encryptSecret("abc", "password", maxCharacters+1); err == nil {
		t.Fatal("too large accepted")
	}
	if _, err := encryptSecret("", "password", 200); err == nil {
		t.Fatal("empty secret accepted")
	}
	ciphertext, err := encryptSecret("private", "passphrase", 90)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := decryptSecret("ET0"+ciphertext[3:], "passphrase"); err == nil {
		t.Fatal("wrong format accepted")
	}
	corrupted := []byte(ciphertext)
	corrupted[20] = '!' // corruption in base64-encoded authenticated payload
	if _, err := decryptSecret(string(corrupted), "passphrase"); err == nil {
		t.Fatal("tampered file accepted")
	}
}

func TestPortableNames(t *testing.T) {
	good := map[string]string{"message": "message.txt", "My file.TXT": "My file.TXT", "test_01": "test_01.txt"}
	for input, expected := range good {
		name, err := validFileName(input)
		if err != nil || name != expected {
			t.Errorf("%q => %q, %v", input, name, err)
		}
	}
	bad := []string{"../evil.txt", "a/b.txt", "C:\\file.txt", "CON.txt", "con.extra.txt", "COM1.txt", "LPT9", "bad.exe", ".hidden", "last..txt", "", strings.Repeat("a", 101) + ".txt"}
	for _, input := range bad {
		if name, err := validFileName(input); err == nil {
			t.Errorf("unsafe name accepted: %q -> %q", input, name)
		}
	}
}

func TestNoOverwriteAndSymlink(t *testing.T) {
	dir := t.TempDir()
	if err := writeEncryptedFile(dir, "file.txt", "ET1-data"); err != nil {
		t.Fatal(err)
	}
	if err := writeEncryptedFile(dir, "file.txt", "overwrite"); !errors.Is(err, os.ErrExist) {
		t.Fatalf("overwrote existing file: %v", err)
	}
	contents, err := readEncryptedFile(dir, "file.txt")
	if err != nil || contents != "ET1-data" {
		t.Fatalf("read failed: %q %v", contents, err)
	}
	if err := os.Symlink(filepath.Join(dir, "file.txt"), filepath.Join(dir, "link.txt")); err == nil {
		if _, err := readEncryptedFile(dir, "link.txt"); err == nil {
			t.Fatal("symlink accepted")
		}
	}
}

func TestStartupBannerOnce(t *testing.T) {
	out := new(bytes.Buffer)
	app := cli{in: bufio.NewReader(strings.NewReader("unexpected\n0\n")), out: out, dir: t.TempDir()}
	if err := app.run(); err != nil && !errors.Is(err, io.EOF) {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out.String(), banner) || strings.Count(out.String(), "by christofrancis") != 1 || !strings.Contains(out.String(), "Goodbye!") {
		t.Fatalf("unexpected menu output: %s", out.String())
	}
}
