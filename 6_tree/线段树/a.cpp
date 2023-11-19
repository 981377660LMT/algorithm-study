// セグ木状に区間を分割したときの処理
//
// 2*offset 個の頂点を持つグラフを考える
// offset+i が元のグラフの頂点 i に対応する
struct DivideInterval {
  int N, offset;
  DivideInterval(int n) : N(n), offset(1) {
    while (offset < N) offset *= 2;
  }
  // 初期化

  // O(N) 根から葉方向へ push する
  // f(p, c) : p -> c へ伝播
  template <typename F>
  void push(const F& f) {
    for (int p = 1; p < offset; p++) {
      f(p, p * 2 + 0);
      f(p, p * 2 + 1);
    }
  }
  // O(N) 葉から根の方向に update する
  // f(p, c1, c2) : c1 と c2 の結果を p へマージ
  template <typename F>
  void update(const F& f) {
    for (int p = offset - 1; p > 0; p--) {
      f(p, p * 2 + 0, p * 2 + 1);
    }
  }

  // [l, r) に対応する index の列を返す
  // 順番は左から右へ並んでいる
  // 例: [1, 11) : [1, 2), [2, 4), [4, 8), [8, 10), [10, 11)
  vector<int> range(int l, int r) {
    assert(0 <= l and l <= r and r <= N);
    vector<int> L, R;
    for (l += offset, r += offset; l < r; l >>= 1, r >>= 1) {
      if (l & 1) L.push_back(l), l++;
      if (r & 1) r--, R.push_back(r);
    }
    for (int i = (int)R.size() - 1; i >= 0; i--) L.push_back(R[i]);
    return L;
  }
  // [l, r) に対応する index に対してクエリを投げる(区間は昇順)
  // f(i) : 区間 i にクエリを投げる
  template <typename F>
  void apply(int l, int r, const F& f) {
    assert(0 <= l and l <= r and r <= N);
    for (int i : range(l, r)) f(i);
  }
};