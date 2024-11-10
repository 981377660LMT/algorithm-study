// LINE(a,b,c): y=(ax+b)/c, 評価点は整数
// 1 次関数の min を [L,R,a,b,c] の列として出力
// c>0, (ax+b)c がオーバーフローしない,
template <typename T>
vc<tuple<T, T, T, T, T>> line_min_function_rational(vc<tuple<T, T, T>> LINE, T L, T R) {
  // 傾き降順
  sort(all(LINE), [&](auto& L, auto& R) -> bool {
    auto& [a1, b1, c1] = L;
    auto& [a2, b2, c2] = R;
    return a1 * c2 > a2 * c1;
  });
  vc<tuple<T, T, T, T, T>> ANS;
  for (auto& [a2, b2, c2]: LINE) {
    while (1) {
      if (ANS.empty()) {
        ANS.eb(L, R, a2, b2, c2);
        break;
      }
      auto& [L1, R1, a1, b1, c1] = ANS.back();
      if ((a1 * L1 + b1) * c2 > (a2 * L1 + b2) * c1) {
        ANS.pop_back();
        if (len(ANS)) get<1>(ANS.back()) = R;
        continue;
      }
      T s = c2 * a1 - a2 * c1;
      if (s == 0) break;
      assert(s > 0);
      T t = b2 * c1 - b1 * c2;
      T x = t / s;
      assert(L1 <= x);
      if (x >= R1 - 1) break;
      R1 = x + 1;
      ANS.eb(x + 1, R, a2, b2, c2);
    }
  }
  return ANS;
}

// LINE(a,b,c): y=(ax+b)/c, 評価点は整数
// 1 次関数の max を [L,R,a,b,c] の列として出力
// c>0, (ax+b)c がオーバーフローしない,
template <typename T>
vc<tuple<T, T, T, T, T>> line_max_function_rational(vc<tuple<T, T, T>> LINE, T L, T R) {
  for (auto& [a, b, c]: LINE) a = -a, b = -b;
  auto ANS = line_min_function_rational<T>(LINE, L, R);
  for (auto& [L, R, a, b, c]: ANS) a = -a, b = -b;
  return ANS;
}

// LINE(a,b): y=ax+b, 評価点は整数
// 1 次関数の min を [L,R,a,b] の列として出力
// ax+b がオーバーフローしない,
template <typename T>
vc<tuple<T, T, T, T>> line_min_function_integer(vc<pair<T, T>> LINE, T L, T R) {
  // 傾き降順
  sort(all(LINE), [&](auto& L, auto& R) -> bool {
    auto& [a1, b1] = L;
    auto& [a2, b2] = R;
    return a1 > a2;
  });
  vc<tuple<T, T, T, T>> ANS;
  for (auto& [a2, b2]: LINE) {
    while (1) {
      if (ANS.empty()) {
        ANS.eb(L, R, a2, b2);
        break;
      }
      auto& [L1, R1, a1, b1] = ANS.back();
      if ((a1 * L1 + b1) > (a2 * L1 + b2)) {
        ANS.pop_back();
        if (len(ANS)) get<1>(ANS.back()) = R;
        continue;
      }
      T s = a1 - a2;
      if (s == 0) break;
      assert(s > 0);
      T t = b2 - b1;
      T x = t / s;
      assert(L1 <= x);
      if (x >= R1 - 1) break;
      R1 = x + 1;
      ANS.eb(x + 1, R, a2, b2);
    }
  }
  return ANS;
}


// LINE(a,b): y=ax+b, 評価点は整数
// 1 次関数の max を [L,R,a,b] の列として出力
// ax+b がオーバーフローしない,
template <typename T>
vc<tuple<T, T, T, T>> line_max_function_integer(vc<pair<T, T>> LINE, T L, T R) {
  for (auto& [a, b]: LINE) a = -a, b = -b;
  auto ANS = line_min_function_integer<T>(LINE, L, R);
  for (auto& [L, R, a, b]: ANS) a = -a, b = -b;
  return ANS;
}