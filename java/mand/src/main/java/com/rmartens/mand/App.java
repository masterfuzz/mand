package com.rmartens.mand;

import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;

import javax.imageio.ImageIO;

import org.apache.commons.math3.complex.Complex;

/**
 * Hello world!
 *
 */
public class App 
{
    public static void main( String[] args )
    {
        System.out.println("Generating...");
        Mandelbrot mand = new Mandelbrot() {{
           center = new Complex(-0.748,0.1);
           zoom = 714.286;
           width = 100;
           height = 100;
           maxIter = 500;
        }};

        BufferedImage img = new BufferedImage(200,200, BufferedImage.TYPE_BYTE_GRAY);
        
        mand.genPixels()
            .forEach(p -> {
                System.out.println(String.format("(%d,%d) -> %s -> %d",
                    p.x, p.y, 
                    mand.pixelToCoordinate(p),
                    mand.f(p).iter
                ));
                img.setRGB(p.x, p.y, p.iter);
            });

        System.out.println("Writing image");
        try {
            ImageIO.write(img, "png", new File("./test.png"));
        } catch (IOException e) {
            e.printStackTrace();
        }
        
        System.out.println("Done");
    }
}
