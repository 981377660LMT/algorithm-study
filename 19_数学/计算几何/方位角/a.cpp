#include "geo/base.hpp"
#include "geo/angle_sort.hpp"

// ベクトルの列が与えられる. 部分列を選んで，和の norm を最小化する.
// 総和の座標の 2 乗和が SM でオーバーフローしないように注意せよ．
// https://atcoder.jp/contests/abc139/tasks/abc139_f
// https://codeforces.com/contest/1841/problem/F
template <typename SM, typename T>
pair<SM, vc<int>> max_norm_sum(vc<Point<T>> dat) {
  auto I = angle_argsort(dat);
  {
    vc<int> J;
    for (auto&& i: I) {
      if (dat[i].x != 0 || dat[i].y != 0) J.eb(i);
    }
    swap(I, J);
  }
  dat = rearrange(dat, I);
  const int N = len(dat);

  if (N == 0) { return {0, {}}; }
  SM ANS = 0;
  pair<int, int> LR = {0, 0};

  int L = 0, R = 1;
  Point<T> c = dat[0];
  auto eval = [&]() -> SM { return SM(c.x) * c.x + SM(c.y) * c.y; };
  if (chmax(ANS, eval())) LR = {L, R};

  while (L < N) {
    Point<T>&A = dat[L], &B = dat[R % N];
    if (R - L < N && (A.det(B) > 0 || (A.det(B) == 0 && A.dot(B) > 0))) {
      c = c + B;
      R++;
      if (chmax(ANS, eval())) LR = {L, R};
    } else {
      c = c - A;
      L++;
      if (chmax(ANS, eval())) LR = {L, R};
    }
  }
  vc<int> ids;
  FOR(i, LR.fi, LR.se) { ids.eb(I[i % N]); }
  return {ANS, ids};
}