package com.rmartens.mand;

import java.util.stream.Stream;
import java.util.stream.IntStream;

import org.apache.commons.math3.complex.Complex;

public class Mandelbrot {
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

    public int f(Complex c) {
        Complex z = new Complex(0,0);

        int i;
        for (i = 0; i < maxIter; i++) {
            if (c.getReal()*c.getReal() + c.getImaginary()*c.getImaginary() > 4) {
                break;
            }
            z = z.pow(2).add(c);
        }
        return i;

    }

    public IntStream generate() {
        return this.getIndecies()
                   .parallel()
                   .map((i) -> this.f(this.indexToCoordinate(i)));
    }

    public Stream<Complex> getCoords() {
        return this.getIndecies()
                   .mapToObj((i) -> this.indexToCoordinate(i));
    }

    public IntStream getIndecies() {
        return IntStream.range(0, this.width * this.height);
    }


}