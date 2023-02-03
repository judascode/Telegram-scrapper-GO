package prompt

import (
	"github.com/manifoldco/promptui"
)

func NameOfGroup() (string, error) {
	promptGroupName := promptui.Prompt{
		Label: "Enter the group name",
	}
	groupName, err := promptGroupName.Run()
	if err != nil {
		return "", err

	}
	return groupName, nil
}
