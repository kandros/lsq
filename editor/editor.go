package editor

import "os"

func Select(editor string) string {
	if editor == "" {
		return "vim"
	}

	checkEnv := os.Getenv(editor)
	if checkEnv != "" {
		return checkEnv
	}

	// Return whatever if not Env Var
	return editor
}
