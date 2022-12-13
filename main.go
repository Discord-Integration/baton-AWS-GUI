package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create a new Fyne app
	a := app.New()

	// Create a new window
	w := a.NewWindow("AWS User Search")
	w.Resize(fyne.NewSize(600, 200))

	input := widget.NewEntry()
	input.SetPlaceHolder("Input user...")

	result := widget.NewLabel("")

	searchButton := widget.NewButton("Search User", func() {

		// Get the text from the entry
		user := input.Text
		fmt.Println("Searching for " + user)

		// Set the label to "Searching..."
		result.SetText("Searching...")

		// When the button is clicked, run the command
		// NOTE: pass user as a parameter to the function and break when found
		displayNames := awsBaton()

		// Print the output of the command to the console
		for i := 0; i < len(displayNames); i++ {

			// If the user is found, set the label to the user
			if user == displayNames[i] {
				response := "User found: " + displayNames[i]
				fmt.Println(response)
				result.SetText(response)
				// Break out of the loop if no more searching is needed
				return
			}
		}

		// If the user is not found, set the label to "User not found"
		notFound := "User not found: " + user
		result.SetText(notFound)
		fmt.Println(notFound)

	})

	clearButton := widget.NewButton("Clear", func() {
		input.SetText("")
		result.SetText("")
	})

	content := container.NewVBox(input, searchButton, clearButton, result)

	// Set the entry as the window's content
	w.SetContent(content)

	// Show the window
	w.ShowAndRun()
}

type AWSResources struct {
	Resources []struct {
		Resource struct {
			DisplayName string `json:"displayName"`
		} `json:"resource"`
	} `json:"resources"`
}

func awsBaton() []string {

	_, err := exec.Command("baton-aws").Output()

	if err != nil {
		log.Fatal(err)
	}

	out, err := exec.Command("baton", "resources", "-o", "json").Output()

	if err != nil {
		log.Fatal(err)
	}

	response := string(out)

	// Unmarshal the byte slice into a AWSResources struct
	var resources AWSResources
	err = json.Unmarshal([]byte(response), &resources)
	if err != nil {
		fmt.Println(err)
	}

	var displayNames []string

	// Loop through the resources and print the displayName for each resource
	for _, r := range resources.Resources {

		displayNames = append(displayNames, r.Resource.DisplayName)
	}

	return displayNames
}
