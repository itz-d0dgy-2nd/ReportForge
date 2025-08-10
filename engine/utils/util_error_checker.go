package Utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func ErrorChecker(_ErrAny error) {

	if _ErrAny != nil {
		if errors.Is(_ErrAny, fs.ErrNotExist) {
			fmt.Printf("::Log:: %s ", _ErrAny.Error())
		} else {
			fmt.Printf("::Error:: An error occurred: %s \n", _ErrAny.Error())
			os.Exit(1)
		}
	}

}
