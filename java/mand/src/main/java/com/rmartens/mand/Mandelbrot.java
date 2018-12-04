package com.rmartens.mand;

import java.util.stream.Stream;
import java.util.stream.IntStream;

import org.apache.commons.math3.complex.Complex;

public class Mandelbrot {
    public static class Pixel {
        public Pixel(int x, int y) {
            this.x = x;
            this.y = y;
            this.iter = 0;
        }
        public Pixel(int x, int y, int iter) {
            this.x = x;
            this.y = y;
            this.iter = iter;
        }
        int x;
        int y;
        int iter;

        @Override
        public String toString() {
            return String.format("(%d,%d)=%d", x, y, iter);
        }
    }
    public Mandelbrot() { }

    public Complex center;
    public double zoom;
    public int maxIter;
    public int width;
    public int height;

    public Complex indexToCoordinate(int i) {
        double cells = (double)this.width * 2;
        double cscale = zoom * cells;
        double x = i / cells;
        double y = i % cells;

        return new Complex(
            (x - cells) / cscale + center.getReal(),
            (y - cells) / cscale - center.getImaginary()
        );
    }

    public Complex pixelToCoordinate(Pixel p) {
        double x = (double)p.x;
        double y = (double)p.y;

        // x in (0,width) -> (-1,1)*scale
        // x' = 2x / width - 1
        return new Complex(
            (2*x / width - 1) / zoom,
            (2*y / height - 1) / zoom
        ).add(center.conjugate());
    }

    // public Pixel indexToPixel(int i) {
    //     int cells = this.width * 2;
    //     return new Pixel(
    //         i / cells,
    //         i % cells
    //     );
    // }

    public int f(Complex c) {
        Complex z = new Complex(0,0);

        int i;
        for (i = 0; i < maxIter; i++) {
            if (z.getReal()*z.getReal() + z.getImaginary()*z.getImaginary() >= 4) {
                return i;
            }
            z = z.pow(2).add(c);
        }
        return 0;

    }

    public Pixel f(Pixel c) {
        return new Pixel(
            c.x, c.y,
            f(pixelToCoordinate(c))
        );
    }

    public IntStream generate() {
        return this.getIndecies()
                   .map((i) -> this.f(this.indexToCoordinate(i)));
    }

    public Stream<Pixel> genPixels() {
        return IntStream.range(0,width)
            .mapToObj(i -> i)
            .flatMap(x -> IntStream.range(0, height).mapToObj(y -> new Pixel(x,y)))
            .map(p -> f(p));
    }

    public Stream<Byte> generateBytes() {
        return this.generate().mapToObj(
            (i) -> new Byte((byte)(i % this.maxIter))
        );
    }

    public Stream<Complex> getCoords() {
        return this.getIndecies()
                   .mapToObj((i) -> this.indexToCoordinate(i));
    }

    public IntStream getIndecies() {
        return IntStream.range(0, this.width * this.height);
    }


}