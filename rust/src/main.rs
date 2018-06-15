extern crate image;
extern crate num_complex;
use num_complex::Complex;

const MAX_ITER: u32 = 2000;
const CENTER_R: f64 = -0.748;
const CENTER_I: f64 = 0.1;
const SCALE: f64 = 714.286;
const CELLS: u32 = 500;
const L255: f64 = 367.88723542668566888;
//const N_THREADS: u32 = 1;


fn main() {
  let mut img = image::GrayImage::new(CELLS*2, CELLS*2);
  let cscale = SCALE * CELLS as f64;

  println!("Running...");
  for (x, y, pixel) in img.enumerate_pixels_mut() {
    let e = escape(Complex::new(
      (x as f64 - CELLS as f64) / cscale + CENTER_R,
      (y as f64 - CELLS as f64) / cscale - CENTER_I
    ));

    //*pixel = image::Luma([(e % 255) as u8]);
    *pixel = image::Luma([f64::floor(L255 * f64::log2(1f64 + e as f64 / MAX_ITER as f64)) as u8]);
  }

  println!("Saving");
  img.save("out.png").unwrap();

  println!("Done");
}

fn escape(c: Complex<f64>) -> u32 {
  let mut z = Complex::new(0f64,0f64);

  for iter in 0..MAX_ITER {
    if z.re * z.re + z.im * z.im  >= 4.0 {
      return iter;
    }

    z = z * z + c
  }
  return 0u32;

}

