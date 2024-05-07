package editor

import (
	"os"
	"os/exec"
)

func File(f string) error {
	editor, ok := os.LookupEnv("VISUAL")
	if !ok {
		editor, ok = os.LookupEnv("EDITOR")
		if !ok {
			editor = "vim"
		}
	}

	c := exec.Command("sh", "-c", editor+" '"+f+"'")
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	return c.Run()
}

func String(src string) (string, error) {
	f, err := os.CreateTemp(os.TempDir(), "jit-edit-*")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write([]byte(src))
	if err != nil {
		return "", err
	}

	filePath := f.Name()

	err = File(filePath)
	if err != nil {
		return "", err
	}

	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
