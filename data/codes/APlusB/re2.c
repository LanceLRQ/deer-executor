#include <stdio.h>
#include <string.h>

int main(int argc, char **argv)
{
	int a, b, p[12][12];
	memset(p, -1, 200 * sizeof(p));
	while (~scanf("%d%d", &a, &b)) {
	    printf("%d\n", a+b);
	}
}
