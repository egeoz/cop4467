package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	hibp "github.com/mattevans/pwned-passwords"
)

var a = app.New()
var w = a.NewWindow("COP4467 Final Project - Password Strength Checker")

var prevChar rune
var score float64
var charCount int
var uniqueCharCount int
var hasDigits bool
var hasSymbols bool
var hasLowerCase bool
var hasUpperCase bool
var isPwned bool

var checkPwned = hibp.NewClient(1)

const windowWidth = 600
const windowHeight = 200

var passwdBox = widget.NewEntry()
var passwordScore = widget.NewProgressBar()

var lblCharCountDefaultText = "Length: "
var lblUniqueCountDefaultText = "Unique count: "
var lblHasLowerCaseDefaultText = "Lower-case: "
var lblHasNumsDefaultText = "Numbers: "
var lblHasUpperCaseDefaultText = "Upper-case: "
var lblHasSpecialDefaultText = "Special characters: "
var lblIsCompromisedDefaultText = "Compromised: "
var lblCharCount = widget.NewLabel(lblCharCountDefaultText)
var lblUniqueCount = widget.NewLabel(lblUniqueCountDefaultText)
var lblHasSpecial = widget.NewLabel(lblHasSpecialDefaultText)
var lblHasNums = widget.NewLabel(lblHasNumsDefaultText)
var lblHasUpperCase = widget.NewLabel(lblHasUpperCaseDefaultText)
var lblHasLowerCase = widget.NewLabel(lblHasLowerCaseDefaultText)
var lblIsCompromised = widget.NewLabelWithStyle(lblIsCompromisedDefaultText, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

const (
	check   = "✓"
	warning = "⚠"
	cross   = '⨯'
)

const (
	digits              = "0123456789"
	symbols             = "!@$&*\"#%'()+/:;<=>?[\\]^{|}~_-.,` "
	lowerCaseCharacters = "abcdefghijklmnopqrstuvwxyz"
	upperCaseCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func checkPasswordStrenght(passwd string) {
	charCount = len(passwd)
	for _, a := range digits {
		if strings.ContainsRune(passwd, a) {
			uniqueCharCount += 1
			hasDigits = true
			if prevChar == a {
				score -= 3
			}

			prevChar = a
		}
	}
	for _, a := range symbols {
		if strings.ContainsRune(passwd, a) {
			uniqueCharCount += 1
			hasSymbols = true
			if prevChar == a {
				score -= 3
			}

			prevChar = a
		}
	}
	for _, a := range lowerCaseCharacters {
		if strings.ContainsRune(passwd, a) {
			uniqueCharCount += 1
			hasLowerCase = true
			if prevChar == a {
				score -= 3
			}

			prevChar = a
		}
	}
	for _, a := range upperCaseCharacters {
		if strings.ContainsRune(passwd, a) {
			uniqueCharCount += 1
			hasUpperCase = true
			if prevChar == a {
				score -= 3
			}

			prevChar = a
		}
	}

}

func btnCheck_onClick() {
	isPwned = false
	hasDigits = false
	hasSymbols = false
	hasUpperCase = false
	hasLowerCase = false
	score = 0
	charCount = 0
	uniqueCharCount = 0
	prevChar = ' '

	checkPasswordStrenght(passwdBox.Text)

	lblCharCount.SetText(fmt.Sprintf(lblCharCountDefaultText+"%d", charCount))
	lblUniqueCount.SetText(fmt.Sprintf(lblUniqueCountDefaultText+"%d", uniqueCharCount))
	lblCharCount.SetText(fmt.Sprintf(lblCharCountDefaultText+"%d", charCount))

	if hasDigits {
		lblHasNums.SetText(fmt.Sprintf(lblHasNumsDefaultText+"%s", "True"))
		score += 10

	} else {
		lblHasNums.SetText(fmt.Sprintf(lblHasNumsDefaultText+"%s", "False"))
		score -= 10
	}

	if hasSymbols {
		lblHasSpecial.SetText(fmt.Sprintf(lblHasSpecialDefaultText+"%s", "True"))
		score += 10

	} else {
		lblHasSpecial.SetText(fmt.Sprintf(lblHasSpecialDefaultText+"%s", "False"))
		score -= 10

	}

	if hasUpperCase {
		lblHasUpperCase.SetText(fmt.Sprintf(lblHasUpperCaseDefaultText+"%s", "True"))
		score += 10

	} else {
		lblHasUpperCase.SetText(fmt.Sprintf(lblHasUpperCaseDefaultText+"%s", "False"))
		score -= 10

	}

	if hasLowerCase {
		lblHasLowerCase.SetText(fmt.Sprintf(lblHasLowerCaseDefaultText+"%s", "True"))
		score += 10

	} else {
		lblHasLowerCase.SetText(fmt.Sprintf(lblHasLowerCaseDefaultText+"%s", "False"))
		score -= 10

	}

	if charCount < 8 || uniqueCharCount < charCount/2 {
		score = 0
	} else {
		score += float64(charCount) * 1.5
		score += float64(uniqueCharCount) * 2
	}

	ret, err := checkPwned.Pwned.Compromised(passwdBox.Text)

	if ret && err == nil {
		lblIsCompromised.SetText(fmt.Sprintf(lblIsCompromisedDefaultText+"%s", "True"))
		score = 0
	} else if err != nil {
		lblIsCompromised.SetText(fmt.Sprintf(lblIsCompromisedDefaultText+"%s", "Error"))
	} else {
		lblIsCompromised.SetText(fmt.Sprintf(lblIsCompromisedDefaultText+"%s", "False"))
	}

	passwordScore.SetValue(score / 100.0)
}

func main() {

	lblTitle := widget.NewLabelWithStyle("Password Strength Checker", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	containerDetails := container.New(layout.NewGridLayout(3), lblCharCount, lblUniqueCount, lblHasNums, lblHasUpperCase, lblHasLowerCase, lblHasSpecial)

	btnCheck := widget.NewButtonWithIcon("Check Password", theme.ConfirmIcon(), btnCheck_onClick)

	containerWindow := container.New(layout.NewGridLayout(1), lblTitle, passwdBox, btnCheck, passwordScore, containerDetails, lblIsCompromised)

	w.SetContent(container.NewMax(containerWindow))
	w.Resize(fyne.NewSize(windowWidth, windowHeight))
	w.ShowAndRun()
}
