#include <stdio.h>
#include <string.h>
#define N 1008016
#define M 6000004

int a[1004][1004];
int b[1004][1004];
int g[M];

int main(int argc, char **argv)
{
    memset(a, -1, N * sizeof(a[0][0]));
    memset(b, -1, N * sizeof(b[0][0]));
    memset(g, -1, M * sizeof(g[0]));
    for(int i = 0;i < M; i++) {
        g[i] = i;
    }
}
