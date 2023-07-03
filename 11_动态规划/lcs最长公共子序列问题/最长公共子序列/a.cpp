#include <iostream>
#include <vector>
#include <map>
#include <set>
#include <algorithm>
#include <stdio.h>
#include <stack>
#include <queue>
#include <atcoder/modint.hpp>
using namespace std;
using namespace atcoder;

typedef unsigned long long ULL;

char s[3001];
char t[3001];

const int N = 3002;
int n, m, u;

struct bitset {
    ULL t[N / 64 + 5];

    bitset() {
        memset(t, 0, sizeof(t));
    }
    bitset(const bitset& rhs) {
        memcpy(t, rhs.t, sizeof(t));
    }

    bool get(int p) {
        return (t[p >> 6] & (1llu << (p & 63))) != 0;
    }
    bitset& set(int p) {
        t[p >> 6] |= 1llu << (p & 63);
        return *this;
    }
    bitset& shift() {
        ULL last = 0llu;
        for (int i = 0; i < u; i++) {
            ULL cur = t[i] >> 63;
            (t[i] <<= 1) |= last, last = cur;
        }
        return *this;
    }

    bitset& operator = (const bitset& rhs) {
        memcpy(t, rhs.t, sizeof(t));
        return *this;
    }
    bitset& operator &= (const bitset& rhs) {
        for (int i = 0; i < u; i++) t[i] &= rhs.t[i];
        return *this;
    }
    bitset& operator |= (const bitset& rhs) {
        for (int i = 0; i < u; i++) t[i] |= rhs.t[i];
        return *this;
    }
    bitset& operator ^= (const bitset& rhs) {
        for (int i = 0; i < u; i++) t[i] ^= rhs.t[i];
        return *this;
    }

    friend bitset operator - (const bitset& lhs, const bitset& rhs) {
        ULL last = 0llu; bitset ret;
        for (int i = 0; i < u; i++) {
            ULL cur = (lhs.t[i] < rhs.t[i] + last);
            ret.t[i] = lhs.t[i] - rhs.t[i] - last;
            last = cur;
        }
        return ret;
    }
} p[26], f[3001], g;


int main() {
    scanf("%s", s);
    scanf("%s", t);
    n = strlen(s);
    m = strlen(t);
    u = n / 64 + 1;
    for (int i = 0, c; i < n; i++) {
        p[s[i] - 'a'].set(i + 1);
    }
    for (int i = 1; i <= m; i++) {
        f[i] = f[i - 1];
        g = f[i];
        g |= p[t[i - 1] - 'a'];
        f[i].shift();
        f[i].set(0);
        f[i] = g - f[i];
        f[i] ^= g;
        f[i] &= g;
    }

    int i = n, j = m;
    string ret;

    while (i && j) {
        if (s[i - 1] == t[j - 1]) {
            ret.push_back(s[i - 1]);
            i--; j--;
        }
        else if (f[j].get(i) == 0) {
            i--;
        }
        else {
            j--;
        }
    }
    std::reverse(ret.begin(), ret.end());
    printf("%s\n", ret.c_str());

    return 0;
}
