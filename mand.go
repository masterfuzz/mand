package main

import (
  "os"
  "fmt"
  "flag"
  "image"
  "image/png"
  "image/color"
)

var (
  maxIter  uint16  = 500
  numRoute uint16  = 4
  delta    float64 = 0.000001
  centerR  float64 = 0
  centerI  float64 = 0
  scale    float64 = 1
  hcells   uint16  = 500
  pixMap   image.NRGBA
  syncChan chan int
)

func getflags() {
  //flag.IntVar(&maxIter, "iter", 500, "Max iterations")
  //flag.IntVar(&numRoute, "routines", 4, "Number of go routines")
  flag.Float64Var(&delta, "delta", 0.000001, "I dont remember this one")
  flag.Float64Var(&centerR, "real", 0, "Real center")
  flag.Float64Var(&centerI, "imag", 0, "Imaginary center")
  flag.Float64Var(&scale, "scale", 1, "Scale")
  //flag.IntVar(&hcells, "hcells", 500, "H Cells")

  flag.Parse()
}

func main() {
  getflags()
	fmt.Printf("Mandelbrot at (%f, %f)\nScale: %fx\nMax iterations: %d\nCells: %d^2\n",
		centerR, centerI,
		scale, maxIter, hcells*2);

  // allocate
  //pixMap = make([][]int, hcells*2)
  //for i := range pixMap {
  //  pixMap[i] = make([]int, hcells*2)
  //}
  pixMap = *image.NewNRGBA(image.Rect(0,0,int(hcells*2),int(hcells*2)))

  fmt.Printf("test: %d\n", escape(0, 0))
  fmt.Printf("test: %d\n", escape(2, 5))

  scale = float64(hcells) * scale

  // start routines
  syncChan = make(chan int)
  for p := 0; p < int(numRoute); p++ {
    go mainLoop(p)
  }

  // join
  for p := 0; p < int(numRoute); p++ {
    fmt.Println(<-syncChan)
  }

  fmt.Println("Writing test.png")

  // write
  outputFile, err := os.Create("test.png")
  if err != nil {
    fmt.Println("FAIL")
  }

  // Encode takes a writer interface and an image interface
  png.Encode(outputFile, &pixMap)

  outputFile.Close()

  fmt.Println("Done!")
}

func mainLoop(offset int) {
  var gray uint16
  var i, j uint16

  for i = uint16(offset); i < hcells*2; i += numRoute {
    for j = 0; j < hcells*2; j++ {
      gray = escape(
        float64(i - hcells) / scale + centerR,
        float64(j - hcells) / scale - centerI,
      ) * 65535 / maxIter

      pixMap.Set(int(i), int(j), color.NRGBA{
        R: uint8(gray & 255),
        G: uint8(gray << 1 & 255),
        B: 0,
        A: 255,
      })
    }
  }
  syncChan <- offset
}

func escape(re0 float64, im0 float64) uint16 {
  var (
    iter uint16 = 0
    re float64 = 0
    im float64 = 0
    tmp float64
  )

  for iter < maxIter && re * re + im * im < 4 {
    // z = z^2 + c
    tmp = re * re - im * im + re0
    im = 2 * re * im + im0
    re = tmp;

    iter++;
  }

  return iter
}

