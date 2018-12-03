package com.rmartens.mand;

import java.awt.image.BufferedImage;
import org.apache.commons.math3.complex.Complex;

/**
 * Hello world!
 *
 */
public class App 
{
    public static void main( String[] args )
    {
        System.out.println( "Hello World!" );
        Mandelbrot mand = new Mandelbrot() {{
           center = new Complex(0,0);
           zoom = 0.1d;
           width = 100;
           height = 100;
           maxIter = 100;
        }};

        mand.getIndecies()
            .mapToObj((i) -> {
                Complex c = mand.indexToCoordinate(i);
                int r = mand.f(c);

                return i + ": f(" + c + ") -> " + r;
            })
            .forEach(
                System.out::println
            );
    }
}
