#include <stdio.h>

static int c1 = 0;
static int c2= 0;
void calc (int l, int r) {
    printf("                    calc(%d,%d)\n",l,r);
    c1++;
    if ( l >= r ) return;
    for (int i = l; i <= r; i++) {
        c2++;
        printf("i=%d;l=%d;r=%d\n",i,l,r);
    }
    printf("                    c1=%d,c2=%d\n",c1,c2);
    int mid = (l + r) / 2;
    if (l == mid) return;
    calc (l, mid);                           //递归1
    if (mid+1 == r) return;                            //分支一般都会有两次递归
    calc (mid + 1, r);                    //递归2
}
int main(void) {
    calc(1,100);
    printf("c1=%d,c2=%d\n",c1,c2);
}
