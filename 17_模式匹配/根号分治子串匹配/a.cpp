//BZOJ 3483

#include <bits/stdc++.h>
using namespace std;

typedef long long LL;
typedef unsigned long long uLL;
const uLL mod = 1e9+7;
const int maxn = 2010;
int n, m, l[maxn];
string a[maxn];
char ch[2000010];
vector <uLL> L[maxn], R[maxn];
int ans;

void input()
{
    scanf("%d", &n);
    for(int i = 1; i <= n; i++){
        scanf("%s", ch);
        l[i] = strlen(ch);
        for(int j = 0; j < l[i]; j++) a[i] += ch[j];
    }
    uLL last;
    for(int i = 1; i <= n; i++){
        last = 171;
        for(int j = 0; j < l[i]; j++){
            last = last*mod+a[i][j];
            L[i].push_back(last);
        }
    }
    for(int i = 1; i <= n; i++){
        last = 171;
        for(int j = 0; j < l[i]; j++){
            last = last*mod+a[i][l[i]-j-1];
            R[i].push_back(last);
        }
    }
}

void work()
{
    scanf("%d", &m);
    ans = 0;
    while(m--)
    {
        uLL x = 171, y = 171;
        scanf("%s", ch);
        int len1 = strlen(ch);
        for(int j = 0; j < len1; j++) ch[j] = (ch[j] - 'a' + ans)%26+'a';
        for(int j = 0; j < len1; j++) x = x*mod + ch[j];
        scanf("%s", ch);
        int len2 = strlen(ch);
        for(int j = 0; j < len2; j++) ch[j] = (ch[j] - 'a' + ans)%26 + 'a';
        for(int j = 0; j < len2; j++) y = y*mod + ch[len2-j-1];
        ans = 0;
        for(int i = 1; i <= n; i++){
            if(L[i][len1-1]==x && R[i][len2-1] == y) ans++;
        }
        printf("%d\n", ans);
    }
}

void output()
{
    //printf("%d\n", ans);
}

int main()
{
    input();
    work();
    //output();
}
