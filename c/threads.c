#include <stdio.h>
#include <pthread.h>

#define THREADS 5
void *tmain(void *arg);

	pthread_t *pth[THREADS];
	//int n[THREADS];
int main(void)
{
	// create threads
	int i = 0;
/*
	for (i; i < THREADS; i++) {
		printf("%d\n", i);
		pthread_create(pth[i], NULL, tmain, NULL);
	}
*/
	pthread_t *pth;
	pthread_create(pth, NULL, tmain, NULL);
	//return 0;

	i = 0;
	while (i < 100) {
		usleep(1);
		printf("main is running...\n");
		i++;
	}

	printf("main waiting on threads...\n");
	pthread_join(*pth, NULL);
	printf("nope\n");
	return 0;
}

void *tmain(void *arg)
{
	int j = 0; //*((int*)arg);
	int i = 0;

	while (i < 100) {
		usleep(1);
		printf("tmain: %d\n", i);
		i++;
	}
	return NULL;
}
