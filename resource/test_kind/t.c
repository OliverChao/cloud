#include <stdio.h>

int age(n,m)
int n;
int m;
{
    printf("%d\n",m);
    int c;
    if (n==1)c=10;

    else{
        c = age(n-1)+2;
    }
    return c;
}
int main(){
    printf("%d\n",age(5));
}
