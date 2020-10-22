#include <stdio.h>

int main(int argc, char **argv)
{
	int a, b, c = 0;
	while (~scanf("%d%d", &a, &b)) {
	    printf("%d\n", a + b);
	    a = b / c;
	}
}
