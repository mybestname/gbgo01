#include <stdio.h>

int fib (int n) {
    if (n < 2 ) return n;
    return fib(n-1) + fib(n-2); //O(k^n)
}

int main(void) {
    int n = fib(42);
    printf("%d\n",n);
    return 0;
}
