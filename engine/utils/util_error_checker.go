package Utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func ErrorChecker(ErrAny error) {

	if ErrAny != nil {
		if errors.Is(ErrAny, fs.ErrNotExist) {
			fmt.Printf("::Log:: %s ", ErrAny.Error())
		} else {
			fmt.Printf("::Error:: An error occurred: %s \n", ErrAny.Error())
			os.Exit(1)
		}
	}

}
