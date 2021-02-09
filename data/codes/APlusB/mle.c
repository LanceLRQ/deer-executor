#include <stdio.h>
#include <string.h>
#define N 100000000

int arr[N];
int main(int argc, char **argv)
{
	int a, b;
    memset(arr, -1, N * sizeof(arr[0]));
	while (~scanf("%d%d", &a, &b)) {
	    printf("%d\n", a+b);
	}
}
