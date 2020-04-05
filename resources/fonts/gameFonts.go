package fonts

import (
	"io/ioutil"
	"log"
)

func getFont(path string) []byte{
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return fontBytes
}

func GetScoreFont() []byte{
	return getFont("resources/fonts/scoreFont.TTF")
}

func GetKeyControllerFont() []byte{
	return getFont("resources/fonts/ComicNeue-Light.ttf")
}