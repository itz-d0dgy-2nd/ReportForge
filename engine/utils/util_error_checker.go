package Utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func ErrorChecker(_ErrAny error) {

	if _ErrAny != nil {
		if errors.Is(_ErrAny, fs.ErrNotExist) {
			fmt.Printf("::Log:: %s ", _ErrAny.Error())
		} else if strings.Contains(_ErrAny.Error(), "executable file not found in $PATH") {
			fmt.Printf("::Error:: chromium based browser not found in $PATH \n")
			os.Exit(1)
		} else {
			fmt.Printf("::Error:: %s \n", _ErrAny.Error())
			os.Exit(1)
		}
	}
}
