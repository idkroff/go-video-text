package main

import (
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/idkroff/go-video-text/internal/config"
	"github.com/idkroff/go-video-text/internal/generator/imagegen"
)

func main() {
	config := config.MustLoad()

	imageGen, err := imagegen.NewGenerator(config.FontPath, config.FontSize, config.MaxWidth)
	if err != nil {
		log.Fatalf("unable to create image generator: %s", err)
	}

	const lettersTimeout = 1
	input := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed mollis nisl tortor. Nullam euismod massa lobortis libero sollicitudin, et vulputate felis cursus. Etiam auctor, nunc id aliquet dignissim, sapien massa consectetur purus, sit amet mollis metus augue at diam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Mauris mattis mattis dui a condimentum. Sed dignissim, ante nec venenatis vehicula, ex nisl egestas est, sed consequat nunc mi vitae lorem. Proin tempus commodo hendrerit. Nam lectus velit, fringilla vel rutrum ac, gravida vitae nibh. Praesent ut dui justo. Praesent metus velit, aliquet id dictum quis, commodo scelerisque leo."

	rows := imageGen.GetRows(input)

	fmt.Println(len(rows))
	for i, v := range rows {
		fmt.Println(i, v)
	}

	//os.Exit(0)

	w, h := imageGen.CalculateWH(rows)
	fmt.Println(w, h)

	img, err := imageGen.NewStringImage(rows[0], w, h)
	for rowI := 1; rowI < len(rows); rowI++ {
		img, err = imageGen.UpdateStringImage(img, rows[rowI], 0, rowI)
	}
	if err != nil {
		log.Fatalf("unable to generate string image: %s", err)
	}

	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	png.Encode(
		f,
		img,
	)
}
