#include <stdlib.h>
#include <stdio.h>
#include <math.h>
#include <pthread.h>
#include <unistd.h>
//#include <png.h>

#define MAX_ITER	500
#define NTHREADS	8
#define DELTA		0.0000001
//#define PVALUE		2.5

int escape(double re, double im);
int escape_p(double re, double im, double p);
static inline void loadBar(int x, int n, int r, int w, float setratio);

u_int16_t* pixmap;
double center_re = 0;
double center_im = 0;
double scale = 1;
double pvalue = 2;
int efficient_power = 1;
int hcells = 500;
int max_iter = MAX_ITER;

void *main_loop(void *arg)
{
	int *offset = (int*)arg;
	int color = 0;
	int i, j;

	for (i = *offset; i < hcells*2; i += NTHREADS) {
		for (j = 0; j < hcells*2; j++) {
			// the case for p=2 is way faster (optimized)
			if (efficient_power) {
				color = escape(
					(i - hcells) / scale + center_re,
					(j - hcells) / scale - center_im);
			} else {
				color = escape_p(
					(i - hcells) / scale + center_re,
					(j - hcells) / scale - center_im,
					pvalue);
			}

//			if (color == MAX_ITER)
//				points++;

			pixmap[j*(hcells*2) + i] = (max_iter < 65536) ? color * 65535 / max_iter : color % 65536;

//			total += color;

		//	pixmap[] = color;
		}
		//printf(".");
		//loadBar(i, hcells*2, hcells/4, 100, 100 * (float)points / (2 * hcells * i));
	}
	return NULL;
}

int main(int argc, char *argv[])
{

	// cmdline args
	// -x, -y	center
	// -z		zoom/scale
	// -n		max iterations
	// -c		cells
	// -p		power
	int c = 0;
	char *outfile = "out";
	while ((c = getopt(argc, argv, "x:y:z:n:c:p:f:")) != -1) {
		double optarg_v = strtod(optarg, NULL);
		switch (c) {
		case 'x':
			center_re = optarg_v;
			break;
		case 'y':
			center_im = optarg_v;
			break;
		case 'z':
			scale = optarg_v;
			break;
		case 'n':
			max_iter = (int)optarg_v;
			break;
		case 'c':
			hcells = (int)optarg_v;
			break;
		case 'p':
			pvalue = optarg_v;
			if (fabs(pvalue - 2.0) < DELTA) {
				efficient_power = 1;
			} else {
				efficient_power = 0;
			}
			break;
		case 'f':
			outfile = optarg;
			break;
		default:
			fprintf(stderr, "Invalid option '-%c'\n", optopt);
			exit(1);
		}
	}


	printf("Mandelbrot at (%f, %f)\nScale: %fx\nMax iterations: %d\nCells: %d^2\n",
		center_re, center_im,
		scale, max_iter, hcells*2);

	// alocate pixel array
	pixmap = malloc(sizeof(int) * 4 * hcells * hcells);
//	u_int32_t total = 0;
//	int points = 0;
	scale = hcells * scale;

	//setbuf(stdout, NULL);
	// create threads
	printf("Running %d threads\n", NTHREADS);
	pthread_t pth[NTHREADS];
	int offsets[NTHREADS];
	int p;
	for (p=0; p<NTHREADS; p++) {
		offsets[p] = p;
		pthread_create(&pth[p], NULL, main_loop, (void*)&offsets[p]);
	}
	// wait for completion
	for (p=0; p<NTHREADS; p++) {
		pthread_join(pth[p], NULL);
	}

	printf("done\n");
/*
	printf("done\nTotal iterations: %d (%d per point avg)\nPoints in set: %d/%d (%3f%%)\n",
		total,
		total / (4*hcells*hcells),
		points,
		4*hcells*hcells,
		100*(float)points / (4*hcells*hcells));
*/

	printf("Writing to file... ");
	FILE *fp;
	fp = fopen(outfile, "w");

	fwrite(pixmap, sizeof(u_int16_t), hcells*hcells*4, fp);		
	fclose(fp);
	free(pixmap);

	printf("done\n");

	return 0;
}


// the standard escape function
// if the complex number (re0 + I * im0) 'escapes' the disk of radius 2,
// it will return the iteration at which it did so. Otherwise returns MAX_ITER
int escape(double re0, double im0)
{
	int iter = 0;
	double re, im;
	double tmp;
	re = 0; im = 0;
//	re =  .1; im = .7;

	while (iter < max_iter && re * re + im * im < 4) {
		// z = z^2 + c
		tmp = re * re - im * im + re0;
		im = 2 * re * im + im0;
		re = tmp;

		iter++;
	}

	return iter;
}

// A less exciting escape function
// Same rules as escape, except the function is simpler: f[z]=z^2
int escape_c(double re0, double im0)
{
	int iter = 0;
	double re, im, tmp;
	re = re0; im = im0;

	while (iter < max_iter && re * re + im * im < 4) {
		// z = z^2
		tmp = re * re - im * im;
		im = 2 * re * im;
		re = tmp;

		iter++;
	}

	return iter;
}

// generalized escape function for arbitrary real powers
// same rules as escape, but uses f[z] = z^p + c
int escape_p(double re0, double im0, double p)
{
	int iter = 0;
	double re, im;
	double z_mod = 0;
	double z_arg = 0;
	re = 0; im = 0;

	while (iter < max_iter && re * re + im * im < 4) {
		// z = z^p + c
		// z = r^p exp(I * arg * p) + c
		z_mod = pow(re * re + im * im, p / 2);
		z_arg = atan2(im, re) * p;

		re = z_mod * cos(z_arg) + re0;
		im = z_mod * sin(z_arg) + im0;

		iter++;
	}

	return iter;
}

// Process has done i out of n rounds,
// and we want a bar of width w and resolution r.
static inline void loadBar(int i, int n, int r, int w, float setratio)
{
	// Only update r times.
	if ( i % (n/r) != 0 ) return;

	// Calculuate the ratio of complete-to-incomplete.
	float ratio = i/(float)n;
	int   c     = ratio * w;

	// Show the percentage complete.
	printf("%3d%% (%3d%%) [", (int)(ratio*100), (int)(setratio) );

	// Show the load bar.
	int x=0;
	for (x=0; x<c; x++) {
		printf("=");
	}

	for (x=c; x<w; x++) {
		printf(" ");
	}
	printf("]");

	// ANSI Control codes to go back to the
	// previous line and clear it.
	// printf("\n33[F33[J");
	printf("\r"); // Move to the first column
	fflush(stdout);
}

