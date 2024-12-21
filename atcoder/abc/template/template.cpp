#pragma GCC target("arch=skylake-avx512")
#pragma GCC target("avx2")
#pragma GCC optimize("O3")
#pragma GCC target("sse4")
#pragma GCC optimize("unroll-loops")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,tune=native")

#include <bits/stdc++.h>
using namespace std;

int main() {
  int N, Q;
  cin >> N >> Q;
  string S;
  cin >> S;
  vector< int > one(N + 1), two(N + 1);
  vector< int > three;
  for(int i = 0; i < N; i++) {
    one[i + 1] = one[i] + (S[i] == '1');
    two[i + 1] = two[i] + (S[i] == '2');
    if(S[i] == '/') three.emplace_back(i);
  }
  while(Q--) {
    int L, R;
    cin >> L >> R;
    --L;
    int l = lower_bound(three.begin(), three.end(), L) - three.begin();
    int r = lower_bound(three.begin(), three.end(), R) - three.begin();
    int ret = 0;
    for(int i = l; i < r; i++) {
      int j = three[i];
      ret = max(ret, min(one[j] - one[L], two[R] - two[j]) * 2 + 1);
    }
    cout << ret << endl;
  }
}