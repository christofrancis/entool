package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

const banner = ` _____       _____           _
| ____|_ __ |_   _|__   ___ | |
|  _| | '_ \  | |/ _ \ / _ \| |
| |___| | | | | | (_) | (_) | |
|_____|_| |_| |_|\___/ \___/|_|
              by christofrancis`

type cli struct {
	in  *bufio.Reader
	out io.Writer
	dir string
}

func (c *cli) ask(prompt string) (string, error) {
	fmt.Fprint(c.out, prompt)
	line, err := c.in.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	if errors.Is(err, io.EOF) && line == "" {
		return "", io.EOF
	}
	return strings.TrimSpace(line), nil
}

func (c *cli) secret(prompt string) (string, error) {
	fmt.Fprint(c.out, prompt)
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintln(c.out)
		return "", errors.New("secret entry requires an interactive terminal")
	}
	value, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(c.out)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (c *cli) encrypt() error {
	name, err := c.ask("Name for the new text file (e.g. message.txt): ")
	if err != nil {
		return err
	}
	name, err = validFileName(name)
	if err != nil {
		return err
	}
	if _, err = os.Lstat(c.dir + string(os.PathSeparator) + name); err == nil {
		return fmt.Errorf("file already exists: %s; nothing overwritten", name)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	secret, err := c.secret("Secret to encrypt (hidden): ")
	if err != nil {
		return err
	}
	if secret == "" {
		return errors.New("the secret cannot be empty")
	}
	passphrase, err := c.secret("Create an encryption key (hidden): ")
	if err != nil {
		return err
	}
	if passphrase == "" {
		return errors.New("the encryption key cannot be empty")
	}
	confirm, err := c.secret("Confirm the encryption key: ")
	if err != nil {
		return err
	}
	if confirm != passphrase {
		return errors.New("the encryption keys do not match; no file was created")
	}
	length, err := c.ask(fmt.Sprintf("How many characters should %s contain? (e.g. 100): ", name))
	if err != nil {
		return err
	}
	count, err := strconv.Atoi(length)
	if err != nil {
		return errors.New("enter a valid whole number")
	}
	ciphertext, err := encryptSecret(secret, passphrase, count)
	if err != nil {
		return err
	}
	if err := writeEncryptedFile(c.dir, name, ciphertext); err != nil {
		return err
	}
	fmt.Fprintf(c.out, "\nEncrypted file created: %s\nLocation: %s\nFile length: %d characters\n", name, c.dir, len(ciphertext))
	fmt.Fprintln(c.out, "Share the .txt file separately from the encryption key.")
	fmt.Fprintln(c.out, "Keep the key safe: without it, the secret cannot be recovered.")
	return nil
}

func (c *cli) decrypt() error {
	matches := listTextFiles(c.dir)
	if len(matches) > 0 {
		end := len(matches)
		if end > 20 {
			end = 20
		}
		fmt.Fprintf(c.out, "Text files in the EnTool folder: %s\n", strings.Join(matches[:end], ", "))
		if len(matches) > 20 {
			fmt.Fprintf(c.out, "... and %d more.\n", len(matches)-20)
		}
	}
	name, err := c.ask("Name of the encrypted text file (e.g. message.txt): ")
	if err != nil {
		return err
	}
	name, err = validFileName(name)
	if err != nil {
		return err
	}
	contents, err := readEncryptedFile(c.dir, name)
	if err != nil {
		return err
	}
	passphrase, err := c.secret("Enter the decryption key (hidden): ")
	if err != nil {
		return err
	}
	secret, err := decryptSecret(contents, passphrase)
	if err != nil {
		return errors.New("unable to decrypt: wrong key, unsupported format, or damaged file")
	}
	fmt.Fprintf(c.out, "\nDecrypted secret: %s\n", secret)
	return nil
}

func (c *cli) run() error {
	fmt.Fprintln(c.out, banner)
	for {
		fmt.Fprintln(c.out, "\n========= EnTool =========")
		fmt.Fprintln(c.out, "1. Encrypt a secret into a new text file")
		fmt.Fprintln(c.out, "2. Decrypt a secret from a text file")
		fmt.Fprintln(c.out, "0. Exit")
		choice, err := c.ask("Select an option (1/2/0): ")
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(c.out, "\nEnTool closed.")
			return nil
		}
		if err != nil {
			return err
		}
		switch choice {
		case "1":
			err = c.encrypt()
		case "2":
			err = c.decrypt()
		case "0":
			fmt.Fprintln(c.out, "Goodbye!")
			return nil
		default:
			fmt.Fprintln(c.out, "Select 1, 2 or 0.")
		}
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(c.out, "\nEnTool closed.")
			return nil
		}
		if err != nil {
			fmt.Fprintf(c.out, "Error: %v\n", err)
		}
	}
}

func main() {
	dir, err := executableFolder()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to locate executable folder:", err)
		os.Exit(1)
	}
	c := &cli{in: bufio.NewReader(os.Stdin), out: os.Stdout, dir: dir}
	if err := c.run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
