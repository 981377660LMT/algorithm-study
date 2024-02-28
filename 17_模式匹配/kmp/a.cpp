#include<bits/stdc++.h>
using namespace std;
/**
 * @brief 计算st从slen开始部分到tlen的前缀函数，前提是<slen的部分已经求出
 *
 * @param st 字符串
 * @param slen 原始长度
 * @param tlen 目标长度，tlen > slen
 * @param pi 前缀函数
 */
void calPi(int* st, int slen, int tlen, int* pi) {
  int j = slen ? pi[slen - 1] : 0;
  for (int i = max(slen, 1); i < tlen; i++) {
    while (j && st[j] != st[i]) {
      if (pi[j - 1] <= j / 2 || st[i] == st[pi[j - 1]])
        j = pi[j - 1];
      else {
        j = pi[j % (j - pi[j - 1]) + (j - pi[j - 1]) - 1];
      }
    }

    if (j || st[j] == st[i]) pi[i] = ++j;
    else pi[i] = 0;
  }
}

const int N = 2e5 + 5;
int s[N], pi[N];
void solve() {
  int n;
  cin >> n;
  int len = 0;
  pi[0] = 0;
  vector<int> ans;
  int m = n;
  while (n--) {
    char op;
    cin >> op;
    if (op == '+') {
      cin >> s[len];
      len++;
      calPi(s, len - 1, len, pi);
    }
    else {
      --len;
    }
    cout << pi[max(len - 1, 0)] << "\n";
  }
}
