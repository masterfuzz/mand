import argparse
import colorsys

class Mandelbrot:
    def __init__(self, args):
        self.num_threads = 1
        self.hcells = 500
        self.max_iter = 500
        self.center = 0 + 0j
        self.scale = 1
        self.pixels = None

    def render(self):
        self._main_loop(0)

    def set_color(self, i, j, c):
        pass

    def _main_loop(self, offset):
        for i in range(offset, self.hcells*2, self.num_threads):
            for j in range(0, self.hcells*2):
                self.set_color(
                    i, j, 
                    self.escape((complex(i - self.hcells, j - self.hcells) + self.center) / self.scale)
                )

    @staticmethod
    def escape(c, maxIter=500):
        i = 0
        z = 0
        while i < maxIter:
            # z -> z^2 + c
            z = z**2 + c
            if z.real**2 + z.imag**2 > 4:
                return i
            i += 1
        return 0

if __name__ == "__main__":
    Mandelbrot(0).render()

