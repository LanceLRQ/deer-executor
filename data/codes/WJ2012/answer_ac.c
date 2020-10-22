#include <stdio.h>

int main() {
    int m, n;
    scanf("%d%*s", &m);
    for (int i = 1; i <= m; i++) {
        int p = 0;
        for (int j = 0; j < 10; j++) {
            if ((i >> j) & 1) {
                p++;
            }
        }
        printf("%d", p);
        for (int j = 0; j < 10; j++) {
            if ((i >> j) & 1) {
                printf(" %d", j + 1);
            }
        }
        printf("\n");
        fflush(stdout);
    }
    scanf("%d%*s", &n);
    int ret = 0;
    for (int i = 0; i < n; i++) {
        int d = 0;
        scanf("%d", &d);
        ret |= 1 << (d-1);
    }
    printf("%d\n", ret);
    fflush(stdout);
}