package main

import (
  "os"
  "math"
  "sync"
  "fmt"
  "flag"
  "image"
  "image/png"
  "image/color"
)

const (
  l2 float64 = 1.4426950408889634074 // 1/math.Log2(2)
  l255 float64 = 367.88723542668566888 //255 / math.Log2(2)
  l2b24 float64 = 24204406.3231229 // 2^24 / math.Log2(2)
)

var (
  maxIter  int     = 500
  numRoute int     = 4
  centerR  float64 = 0
  centerI  float64 = 0
  scale    float64 = 1
  hcells   int     = 500
  outPath  string  = "test.png"
  colorScheme string = "rgb"
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
  flag.StringVar(&colorScheme, "color", "rgb", "Color scheme")

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

    for i := offset; i < hcells*2; i += numRoute {
      for j := 0; j < hcells*2; j++ {
        pixMap.Set(i, j, getColor(escape(
            (float64(i) - float64(hcells)) / scale + centerR,
            (float64(j) - float64(hcells)) / scale - centerI,
          )))
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

  if iter == maxIter {
    return 0
  }
  return uint32(iter)
}

func getColor(i uint32) color.NRGBA {
  if i == 0 {
    return color.NRGBA{}
  }
  switch colorScheme {
  default:
    fallthrough
  case "rgbweird":
    v := i * uint32(16777216 / maxIter)
    return color.NRGBA{
      R: uint8(v & 255),
      G: uint8(v >> 8 & 255),
      B: uint8(v >> 16 & 255),
      A: 255,
    }
  case "rgb":
    return HSVtoRGB(float64(i) / float64(maxIter), 1, 1)
  case "logrgb":
    v := l2 * math.Log2(1 + float64(i) / float64(maxIter))
    return HSVtoRGB(v, 1, 1)
  case "gray":
    v := uint8(math.Floor(255 * float64(i) / float64(maxIter)))
    return color.NRGBA{
      R: v,
      G: v,
      B: v,
      A: 255,
    }
  case "loggray":
    v := uint8(math.Floor(l255 * math.Log2(1 + float64(i) / float64(maxIter))))
    return color.NRGBA{
      R: v,
      G: v,
      B: v,
      A: 255,
    }

  }
}

func HSVtoRGB(h, s, v float64) color.NRGBA {
  var r, g, b float64
  h = 360 * h
  c := v * s
  x := c * (1 - math.Abs(math.Mod((h / 60), 2) -1))
  m := v - c

  switch {
  case h < 60:
    r, g, b = c, x, 0
  case h < 120:
    r, g, b = x, c, 0
  case h < 180:
    r, g, b = 0, c, x
  case h < 240:
    r, g, b = 0, x, c
  case h < 300:
    r, g, b = x, 0, c
  default:
    r, g, b = c, 0, x
  }

  return color.NRGBA{
    R: uint8(math.Floor((r + m)*255)),
    G: uint8(math.Floor((g + m)*255)),
    B: uint8(math.Floor((b + m)*255)),
    A: 255,
  }
}

