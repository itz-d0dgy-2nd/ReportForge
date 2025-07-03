package Utils

import (
	"fmt"
	"errors"
	"io/fs"
)

func ErrorChecker(ErrAny error) {

	if ErrAny != nil {
		if errors.Is(ErrAny, fs.ErrNotExist){
			fmt.Printf("::Log:: %s", ErrAny.Error())	
		} else {
			fmt.Printf("::Error:: An error occurred: %s", ErrAny.Error())
			panic(ErrAny)
		}
	}

}
