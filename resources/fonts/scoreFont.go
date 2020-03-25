package fonts

import (
	"io/ioutil"
	"log"
)

func GetFont() []byte{
	fontBytes, err := ioutil.ReadFile("resources/fonts/scoreFont.TTF")
	if err != nil {
		log.Fatal(err)
	}
	return fontBytes
}