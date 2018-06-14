
const MAX_ITER: i32 = 500;
const CENTER_R: f64 = 0.0;
const CENTER_I: f64 = 0.0;
const SCALE: f64 = 1.0;
const CELLS: u32 = 500;
const N_THREADS: u32 = 1;


fn main() {
  println!("Hello! {}", escape(1.0, 1.0));
  main_loop(0);
  println!("Done!");
}

fn main_loop(offset: u32) {
  let mut i = offset;

  while i < CELLS*2 {
    for j in 0..(CELLS*2) {
      escape(
        (i as f64 - CELLS as f64) / SCALE + CENTER_R,
        (j as f64 - CELLS as f64) / SCALE - CENTER_I
      );
    }
    i += N_THREADS;
  }
}

fn escape(re0: f64, im0: f64) -> i32 {
  let mut re = 0f64;
  let mut im = 0f64;
  let mut tmp = 0f64;

  for iter in 0..MAX_ITER {
    if re * re + im * im  >= 4.0 {
      return iter;
    }

    tmp = re * re - im * im + re0;
    im = 2.0 * re *im + im0;
    re = tmp;
  }
  return 0;
}
