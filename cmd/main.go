package main

import (
	"flag"
	"fmt"
	"mime"
	"strings"

	"github.com/bclicn/color"

	"golang.org/x/sys/windows/registry"
)

var (
	flagFix = flag.Bool("fix", false, "Fix detected issue(s)")
)

func main() {
	flag.Parse()

	fmt.Println(color.Cyan("[ Diagnosting current state ]"))

	if checkRoutine() {
		return
	}

	if !*flagFix {
		fmt.Println(color.Yellow(
			"\nSkipping fixing routine.\nRun with flag -fix to fix detected issue(s)."))
		return
	}

	fmt.Println(color.Cyan("\n[ Issue fix routine ]"))

	err := fixRegValue()
	if err != nil {
		fmt.Printf(color.Red("\nERROR ► Failed setting registry value: %s"), err.Error())
	}

	fmt.Println(color.Green("\nEverything should be fine now. Re-run to check state."))
}

func checkRoutine() bool {
	mimeVal, mimeValid := checkIfErrorOccurs()
	printStatusLine("mime.TypeByExtension(\".js\")", mimeVal, mimeValid)
	if mimeValid {
		fmt.Println(color.Green("\nEverything is as expected."))
		return true
	}

	regVal, regValid, err := checkCurrRegValue()
	if err != nil {
		fmt.Printf(color.Red("\nERROR ► Failed reading registry value: %s"), err.Error())
		return true
	}
	printStatusLine("reg query HKCR\\.js Content Type", regVal, regValid)
	if regValid {
		fmt.Printf(color.Yellow(
			"mime type was detected falsely even though registry key value of " +
				"HKCR\\.js Content Type is as expected. The issue must be caused " +
				"from some other factor."))
		return true
	}

	return false
}

func isValidContentType(ct string) bool {
	return strings.HasPrefix(ct, "application/javascript")
}

func checkIfErrorOccurs() (string, bool) {
	res := mime.TypeByExtension(".js")
	valid := isValidContentType(res)

	return res, valid
}

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

	valid := isValidContentType(res)

	return res, valid, nil
}

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

func printStatusLine(desc, val string, valid bool) {
	validS := color.BGreen("as expected")
	if !valid {
		validS = color.BRed("not as expected")
	}

	fmt.Printf("%s → '%s' ► %s\n",
		desc, color.Yellow(val), validS)
}
