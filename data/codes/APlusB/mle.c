#include <stdio.h>
#include <string.h>
#define M 10000004

int main(int argc, char **argv)
{
	int g[M];
    memset(g, -1, M * sizeof(g[0]));
    for(int i = 0;i < M; i++) {
        g[i] = i;
    }
}
