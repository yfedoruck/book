package fail

import "log"

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Warning(err error) {
	if err != nil {
		log.Println(err)
	}
}
