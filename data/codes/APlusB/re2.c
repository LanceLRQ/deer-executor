#include <stdio.h>

int main(int argc, char **argv)
{
	int a, b, *p;
	*p = 1;
	while (~scanf("%d%d", &a, &b)) {
	    printf("%d\n", a+b);
	}
}
