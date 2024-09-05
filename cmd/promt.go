package cmd

import "github.com/manifoldco/promptui"

func ask(question string, defaultValues ...bool) bool {
	cursorPos := 0
	if len(defaultValues) > 0 && !defaultValues[0] {
		cursorPos = 1
	}
	prompt := promptui.Select{
		Label: question,
		Items: []string{
			"Yes",
			"No",
		},
		CursorPos: cursorPos,
	}

	_, val, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	return val == "Yes"
}
