package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func main() {
	imgToBeDecoded := flag.String("decode", "", "specify the image to be decoded")
	imgToBeEncoded := flag.String("encode", "", "specify the image to be encoded")
	flag.Parse()

	if *imgToBeDecoded != "" {
		fmt.Println(Decode(*imgToBeDecoded))
	} else if *imgToBeEncoded != "" {
		text := flag.Arg(0)
		if text == "" {
			fmt.Println("text not provided")
			os.Exit(1)
		}
		Encode(*imgToBeEncoded, text)
	} else {
		fmt.Println("arguments not provided")
		os.Exit(1)
	}
}

func Encode(path string, text string) {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	img, err := png.Decode(file)

	if err != nil {
		fmt.Println("Error: PNG could not be decoded")
		os.Exit(1)
	}

	file.Close()

	length := len(text)
	size := img.Bounds()
	width, height := size.Dx(), size.Dy()

	if length > width*height {
		fmt.Println("Image is not large enough to encode the given text")
		os.Exit(1)
	}

	m := image.NewRGBA(image.Rect(0, 0, width, height))

	x, y := 0, 0
	var r, g, b, a uint8

	for i := 0; i < length; i++ {
		bin := IntToBinary(int(text[i]))
		r, g, b, a = ToUint8(img.At(x, y).RGBA())

		if bin[0] == 1 {
			r = MakeOdd(r)
		} else {
			r = MakeEven(r)
		}

		if bin[1] == 1 {
			g = MakeOdd(g)
		} else {
			g = MakeEven(g)
		}

		if bin[2] == 1 {
			b = MakeOdd(b)
		} else {
			b = MakeEven(b)
		}

		m.SetRGBA(x, y, color.RGBA{r, g, b, a})

		if x == width-1 {
			x = 0
			y++
		} else {
			x++
		}

		r, g, b, a = ToUint8(img.At(x, y).RGBA())

		if bin[3] == 1 {
			r = MakeOdd(r)
		} else {
			r = MakeEven(r)
		}

		if bin[4] == 1 {
			g = MakeOdd(g)
		} else {
			g = MakeEven(g)
		}

		if bin[5] == 1 {
			b = MakeOdd(b)
		} else {
			b = MakeEven(b)
		}

		m.SetRGBA(x, y, color.RGBA{r, g, b, a})

		if x == width-1 {
			x = 0
			y++
		} else {
			x++
		}

		r, g, b, a = ToUint8(img.At(x, y).RGBA())

		if bin[6] == 1 {
			r = MakeOdd(r)
		} else {
			r = MakeEven(r)
		}

		if bin[7] == 1 {
			g = MakeOdd(g)
		} else {
			g = MakeEven(g)
		}

		if i == length-1 {
			b = MakeEven(b)
		} else {
			b = MakeOdd(b)
		}

		m.SetRGBA(x, y, color.RGBA{r, g, b, a})

		if x == width-1 {
			x = 0
			y++
		} else {
			x++
		}
	}

	for ; y < height; y++ {
		for ; x < width; x++ {
			r, g, b, a = ToUint8(img.At(x, y).RGBA())
			m.SetRGBA(x, y, color.RGBA{r, g, b, a})
		}
		x = 0
	}

	file, err = os.OpenFile(path, os.O_WRONLY, os.ModeSetuid)

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	err = png.Encode(file, m)

	if err != nil {
		fmt.Println("Error: PNG could not be encoded")
		os.Exit(1)
	}
}

func Decode(path string) string {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	img, err := png.Decode(file)

	if err != nil {
		fmt.Println("Error: PNG could not be decoded")
		os.Exit(1)
	}

	file.Close()

	size := img.Bounds()
	width, height := size.Dx(), size.Dy()

	var bin []int

outerloop:
	for y := 0; y < height; y++ {
		i := 0
		for x := 0; x < width; x++ {
			r, g, b, _ := ToUint8(img.At(x, y).RGBA())
			d1, d2, d3 := 1, 1, 1
			if r%2 == 0 {
				d1 = 0
			}
			if g%2 == 0 {
				d2 = 0
			}
			if i == 2 {
				i = 0
				bin = append(bin, d1, d2)
				if b%2 == 0 {
					break outerloop
				}
				continue
			} else if b%2 == 0 {
				d3 = 0
			}
			bin = append(bin, d1, d2, d3)
			i++
		}
	}
	return BinaryToText(bin)
}

func IntToBinary(n int) []int {
	bin := make([]int, 8)
	for j := 7; j > -1; j-- {
		bin[j] = n % 2
		n = n / 2
	}
	return bin
}

func BinaryToText(bin []int) (str string) {
	for i := 0; i < len(bin)/8; i++ {
		charcode := 0
		for j := 7; j > -1; j-- {
			k := float64(7 - j)
			charcode += bin[i*8+j] * int(math.Pow(2, k))
		}
		str += string(rune(charcode))
	}
	return
}

func MakeOdd(n uint8) uint8 {
	if n == 0 {
		return 1
	} else if n%2 == 1 {
		return n
	}
	return n - 1
}

func MakeEven(n uint8) uint8 {
	if n%2 == 0 {
		return n
	}
	return n - 1
}
func ToUint8(a, b, c, d uint32) (u, v, w, x uint8) {
	u = uint8(a)
	v = uint8(b)
	w = uint8(c)
	x = uint8(d)
	return
}
