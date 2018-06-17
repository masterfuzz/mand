extern crate image;
extern crate num_complex;
extern crate rayon;
use num_complex::Complex;
use rayon::prelude::*;

const MAX_ITER: u32 = 2000;
const CENTER_R: f64 = -0.748;
const CENTER_I: f64 = 0.1;
const SCALE: f64 = 714.286;
const CELLS: u32 = 500;
const L255: f64 = 367.88723542668566888;
//const N_THREADS: u32 = 1;


fn main() {
  let cscale = SCALE * CELLS as f64;
  let mut hist = vec![0u32; MAX_ITER as usize];

  println!("Running...");
  let escapes: Vec<_> = (0..(CELLS*CELLS*4)).into_par_iter().map(|i| {
    let x = i / (CELLS*2);
    let y = i % (CELLS*2);
    let e = escape(Complex::new(
      (x as f64 - CELLS as f64) / cscale + CENTER_R,
      (y as f64 - CELLS as f64) / cscale - CENTER_I
    ));

    return e; //f64::floor(L255 * f64::log2(1f64 + e as f64 / MAX_ITER as f64)) as u8;
  }).collect();

  for v in escapes.iter() {
    hist[*v as usize] += 1;
  }
  let total: u32 = hist.iter().sum();

  let pixels: Vec<_> = escapes.iter().map(|e| {
    if e == 0 {
      return 0u8;
    }
    let hnorm: f64 = (0..*e).map(|j| hist[j as usize] as f64 / total as f64).sum();
    return f64::floor(L255 * f64::log2(1f64 + hnorm)) as u8;
  }).collect();

  println!("Saving");
  let img = image::GrayImage::from_raw(CELLS*2, CELLS*2, pixels);
  img.unwrap().save("out.png").unwrap();

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

