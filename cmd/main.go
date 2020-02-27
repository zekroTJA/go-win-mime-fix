package main

import (
	"flag"
	"fmt"
	"mime"

	"github.com/bclicn/color"
	"github.com/zekroTJA/go-win-mime-fix/internal/helper"

	"golang.org/x/sys/windows/registry"
)

var (
	flagFix = flag.Bool("fix", false, "Fix detected issue(s)")
)

func main() {
	flag.Parse()

	// Check current state
	fmt.Println(color.Cyan("[ Diagnosting current state ]"))
	if checkRoutine() {
		return
	}

	// If flag -fix is not passed, show this warning
	// and stop program.
	if !*flagFix {
		fmt.Println(color.Yellow(
			"\nSkipping fixing routine.\nRun with flag -fix to fix detected issue(s)."))
		return
	}

	// Run fix routine
	fmt.Println(color.Cyan("\n[ Issue fix routine ]"))
	err := fixRegValue()
	if err != nil {
		fmt.Printf(color.Red("\nERROR ► Failed setting registry value: %s"), err.Error())
		return
	}

	// Final success message
	// The program must be re-started manually instead of re-running
	// checkRoutine(), because the values for mime.TypeByExtension
	// will be read from registry and set on runtime startup.
	fmt.Println(color.Green("\nEverything should be fine now. Re-run to check state."))
}

// checkRoutine executes both checks
// checkIfErrorOccurs() and checkCurrRegValue(),
// formats and prints the results. If the returned
// value is true, that means for the main routine,
// that the program should be stopped at this point.
func checkRoutine() bool {
	mimeVal, mimeValid := checkIfErrorOccurs()
	helper.PrintStatusLine("mime.TypeByExtension(\".js\")", mimeVal, mimeValid)
	if mimeValid {
		fmt.Println(color.Green("\nEverything is as expected."))
		return true
	}

	regVal, regValid, err := checkCurrRegValue()
	if err != nil {
		fmt.Printf(color.Red("\nERROR ► Failed reading registry value: %s"), err.Error())
		return true
	}
	helper.PrintStatusLine("reg query HKCR\\.js Content Type", regVal, regValid)
	if regValid {
		fmt.Printf(color.Yellow(
			"mime type was detected falsely even though registry key value of " +
				"HKCR\\.js Content Type is as expected. The issue must be caused " +
				"from some other factor."))
		return true
	}

	return false
}

// checkIfErrorOccurs checks if the return
// value of mime.TypeByExtension(".js") is
// an unexpected value.
// The recovered value will be returned.
// If the recovered value is indeed unexpected,
// the returned bool value will be false.
func checkIfErrorOccurs() (string, bool) {
	res := mime.TypeByExtension(".js")
	valid := helper.IsValidContentType(res)

	return res, valid
}

// checkCurrRegValue checks if the registry value,
// which is suspect to be the cause of the issue,
// is actually set to an unexpected value.
// If this is indeed the case, this function returns
// the actually recovered value and false.
// If reading of the registry key was errous,
// the error will be returned.
func checkCurrRegValue() (string, bool, error) {
	k, err := registry.OpenKey(registry.CLASSES_ROOT, ".js", registry.QUERY_VALUE)
	if err != nil {
		return "", false, err
	}

	defer k.Close()

	res, _, err := k.GetStringValue("Content Type")
	if err != nil {
		return "", false, err
	}

	valid := helper.IsValidContentType(res)

	return res, valid, nil
}

// fixRegValue tries to set the issue causing
// registry key value to the expected value.
// If this fails, the error will be returned.
func fixRegValue() error {
	fmt.Print("Fixing registry key value HKCR\\.js Content Type...")
	k, err := registry.OpenKey(registry.CLASSES_ROOT, ".js", registry.SET_VALUE)
	if err != nil {
		fmt.Println(" → " + color.Red("failed"))
		return err
	}

	defer k.Close()

	err = k.SetStringValue("Content Type", "application/javascript")
	if err != nil {
		fmt.Println(" → " + color.Red("failed"))
		return err
	}

	fmt.Println(" → " + color.Green("success"))
	return nil
}
