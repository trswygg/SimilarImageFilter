package imagefilter

import (
	"fmt"
)

type Style int

type DrawFunc func(SimilarPicture)

func (C *Comparator) Show(needTrim bool, drawFunc DrawFunc) {
	if needTrim {
		C.trimResult('a')
		C.trimResult('d')
		C.trimResult('p')
	}

	fmt.Printf("\n")

	if drawFunc == nil {
		panic("drawFunc can`t be nil")
	}
	for group := C.similarPictures.Front(); group != nil; group = group.Next() {
		g := group.Value.(SimilarPicture)
		drawFunc(g)
	}
}

