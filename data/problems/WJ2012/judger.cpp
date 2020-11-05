/*
@author: Wolf Zheng
@algorithm: check code
*/

#include <iostream>

using namespace std;

bool flag[12];
bool rightflag[12];

int main(int argc, char *argv[])
{
    FILE* fin = fopen(argv[1], "r");
	FILE* fp = fopen(argv[4], "w+");
	//FILE* fp = fopen("test_out.txt", "w+");
	int n, pois;
	fscanf(fin, "%d %d", &n, &pois);
	//scanf("%d %d", &n, &pois);
	printf("%d cans SCAU yogurt\n", n);
	fflush(stdout);
	int ansnum = 0;
	for(int i = 1; i <= n; i++)
    {
        int num;
        char che;
        scanf("%d%c", &num, &che);
        //printf("t%d\n", num);

        if(num > 10 || num < 0)
        {
            fprintf(fp, "wrong mouse num: %d\n", num);
            fclose(fp);
            return 4;
        }
        if(num != 0 && che != ' ')
        {
            fprintf(fp, "wrong format1: %d%c\n", num, che);
            fclose(fp);
            return 4;
        }
        if(num == 0 && che != '\n')
        {
            fprintf(fp, "wrong format2: %d%c\n", num, che);
            fclose(fp);
            return 4;
        }
        fprintf(fp, "yogurt %d eated by %d mice.\n", i, num);

        memset(flag, false, sizeof(flag));
        for(int j = 0; j < num - 1; j++)
        {
            int mice;
            scanf("%d%c", &mice, &che);
            //printf(" t%d ", mice);
            if(che != ' ')
            {
                fprintf(fp, "wrong format3: %d%d\n", mice, che);
                fclose(fp);
                return 4;
            }
            if(mice > 10 || mice < 1)
            {
                fprintf(fp, "wrong mouse index: %d\n", mice, che);
                fclose(fp);
                return 4;
            }
            if(flag[mice])
            {
                fprintf(fp, "mouse %d been used.\n", mice, che);
                fclose(fp);
                return 4;
            }
            fprintf(fp, "%d ", mice);
            flag[mice] = true;
            if(i == pois)
            {
                rightflag[mice] = true;
                ansnum++;
            }
        }
        int mice;
        scanf("%d%c", &mice, &che);
        if(che != '\n')
        {
            fprintf(fp, "wrong format4: %d%c\n", mice, che);
            fclose(fp);
            return 4;
        }
        if(mice > 10 || mice < 1)
        {
            fprintf(fp, "wrong mouse index: %d\n", mice, che);
            fclose(fp);
            return 4;
        }
        if(flag[mice])
        {
            fprintf(fp, "mouse %d been used.\n", mice, che);
            fclose(fp);
            return 4;
        }
        fprintf(fp, "%d\n", mice);
        flag[mice] = true;
        if(i == pois)
        {
            rightflag[mice] = true;
            ansnum++;
        }
    }
    printf("%d mice died\n", ansnum);
    int pnum = 0;
    for(pnum = 1; pnum <= 10; pnum++)
    {
        if(rightflag[pnum])
        {
            printf("%d", pnum);
            break;
        }
    }
    pnum++;
    for(; pnum <= 10; pnum++)
    {
        if(rightflag[pnum])
        {
            printf(" %d", pnum);
        }
    }
    printf("\n");
    fflush(stdout);
    int uans;
    scanf("%d", &uans);
    if(uans == pois)
    {
        fprintf(fp, "Ans accepted.\n");
        fclose(fp);
        return 0;
    }
    else
    {
        fprintf(fp, "Wrong answer. User's ans is %d.\n", uans);
        fclose(fp);
        return 4;
    }
	fclose(fp);
	return 4;
}

/*
1 1
1 2
2 1 2
1 3
2 1 3
2 2 3
3 1 2 3
1 4
2 1 4
2 2 4

1 1
1 2
2 1 2
1 3
2 1 3
2 2 3
3 1 2 3
1 4
2 1 4
2 2 4
3 1 2 4
2 3 4
3 1 3 4
3 2 3 4
4 1 2 3 4
 */
