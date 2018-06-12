package main

import (
  "os"
  "sync"
  "fmt"
  "flag"
  "image"
  "image/png"
  "image/color"
)

var (
  maxIter  int     = 500
  numRoute int     = 4
  centerR  float64 = 0
  centerI  float64 = 0
  scale    float64 = 1
  hcells   int     = 500
  outPath  string  = "test.png"
  pixMap   image.NRGBA
)

func getflags() {
  flag.IntVar(&maxIter, "iter", 500, "Max iterations")
  flag.IntVar(&numRoute, "routines", 4, "Number of go routines")
  flag.Float64Var(&centerR, "real", 0, "Real center")
  flag.Float64Var(&centerI, "imag", 0, "Imaginary center")
  flag.Float64Var(&scale, "scale", 1, "Scale")
  flag.StringVar(&outPath, "out", "test.png", "Path to output PNG")
  flag.IntVar(&hcells, "size", 500, "Number of pixels")

  flag.Parse()
}

func main() {
  getflags()
	fmt.Printf("Mandelbrot at (%f, %f)\nScale: %fx\nMax iterations: %d\nCells: %d^2\n",
		centerR, centerI,
		scale, maxIter, hcells*2);

  // allocate image
  pixMap = *image.NewNRGBA(image.Rect(0,0,hcells*2,hcells*2))

  // adjust scale with cells
  scale = float64(hcells) * scale

  // start routines
  var wg sync.WaitGroup
  for p := 0; p < numRoute; p++ {
    mainLoop(p, &wg)
  }

  // join
  wg.Wait()

  fmt.Printf("Writing %v\n", outPath)

  // write
  outputFile, err := os.Create(outPath)
  if err != nil {
    fmt.Println("FAIL")
  }

  // Encode takes a writer interface and an image interface
  png.Encode(outputFile, &pixMap)

  outputFile.Close()

  fmt.Println("Done!")
}

func mainLoop(offset int, wg *sync.WaitGroup) {
  wg.Add(1)

  go func() {
    defer wg.Done()
    var gray uint32

    for i := offset; i < hcells*2; i += numRoute {
      for j := 0; j < hcells*2; j++ {
        gray = escape(
          (float64(i) - float64(hcells)) / scale + centerR,
          (float64(j) - float64(hcells)) / scale - centerI,
        ) * uint32(16777216 / maxIter)

        pixMap.Set(i, j, color.NRGBA{
          R: uint8(gray & 255),
          G: uint8(gray << 1 & 255),
          B: uint8(gray << 2 & 255),
          A: 255,
        })
      }
    }
  }()
}

func escape(re0 float64, im0 float64) uint32 {
  var (
    iter int = 0
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

  return uint32(iter)
}

