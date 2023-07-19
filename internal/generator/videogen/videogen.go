package videogen

import vidio "github.com/AlexEidt/Vidio"

func NewTest() {
	w, h := 100, 100

	options := vidio.Options{FPS: 1, Loop: 0, Delay: 1000}
	v, _ := vidio.NewVideoWriter("output.mp4", w, h, &options)
	defer v.Close()

}
