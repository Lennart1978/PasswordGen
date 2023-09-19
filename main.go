// generateRandomPassword generates a random password string of the given length.
//
// It uses a cryptographically secure random number generator seeded with
// the current time to generate random runes from the set of alphanumeric
// characters and symbols in allChars.
//
// The runes are assembled into a string and returned.
//
// Parameters:
//   - numDigits: The length of the password to generate
//
// Returns:
//   - A random password string of length numDigits
//func generateRandomPassword(numDigits int) string {
// ...
// }

// SaveStringToFile saves the given string to the specified file path.
//
// It creates or truncates the file and writes the string to it.
//
// Parameters:
//   - filename: The path of the file to write to
//   - data: The string to save
//
// Returns:
//   - An error if one occurred while writing the file, otherwise nil
//func SaveStringToFile(filename, data string) error {
// ...
//}

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/gen2brain/beeep"
)

const (
	allChars     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~"
	windowWidth  = 300
	windowHeight = 150
	defaultFile  = "password.txt"
)

func main() {
	// Create app and main window
	a := app.NewWithID("org.PasswordGen.app")
	w := a.NewWindow("Password Generator")
	w.Resize(fyne.NewSize(windowWidth, windowHeight))

	// Show About window
	w.SetMainMenu(createMenu())

	// Create UI elements
	label1 := widget.NewLabel("Number of digits:")
	entry1 := widget.NewEntry()
	entry2 := widget.NewEntry()
	generateBtn := widget.NewButton("Generate", func() {
		// Generate password
		numDigits, err := strconv.Atoi(entry1.Text)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		password := generateRandomPassword(numDigits)
		entry2.SetText(password)
	})

	copyBtn := widget.NewButton("Copy", func() {
		// Copy password
		w.Clipboard().SetContent(entry2.Text)
	})

	saveBtn := widget.NewButton("Save", func() {
		// Save password to file
		savePasswordToFile(w, entry2.Text)
	})

	// Layout UI
	content := container.NewVBox(
		label1, entry1, entry2,
		generateBtn,
		container.NewHBox(copyBtn, saveBtn))

	w.SetContent(content)

	// Show window
	w.ShowAndRun()
}

func createMenu() *fyne.MainMenu {

	aboutItem := fyne.NewMenuItem("Über", func() {
		err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		if err != nil {
			panic(err)
		}
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "About",
			Content: "2023 - PasswordGenerator 1.0\nby Lennart Martens",
		})
	})

	return fyne.NewMainMenu(
		fyne.NewMenu("Help", aboutItem),
	)
}

// Extracted function to save password
func savePasswordToFile(w fyne.Window, password string) {
	fileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		// Error handling
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		if writer == nil {
			return
		}

		defer writer.Close()

		// Save password
		err = SaveStringToFile(writer.URI().Path(), password)
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Saved", "File saved successfully", w)
		}
	}, w)

	fileDialog.SetFileName(defaultFile)
	fileDialog.Show()
}

func generateRandomPassword(numDigits int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	passwordRunes := make([]rune, numDigits)
	for i := range passwordRunes {
		passwordRunes[i] = rune(allChars[r.Intn(len(allChars))])
	}

	return string(passwordRunes)
}

func SaveStringToFile(filename, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(data); err != nil {
		return err
	}

	fmt.Printf("String saved to %s\n", filename)
	return nil
}
