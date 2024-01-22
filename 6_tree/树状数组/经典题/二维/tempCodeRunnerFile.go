
// LL(N, Q);
//   using T = tuple<int, int, int, int, int>;
//   VEC(T, dat, N);
//   vc<T> QUERY;
//   vc<int> X, Y;
//   FOR(Q) {
//     INT(t);
//     if (t == 0) {
//       LL(a, b, c, d, x);
//       QUERY.eb(a, b, c, d, x);
//     }
//     if (t == 1) {
//       LL(x, y);
//       X.eb(x), Y.eb(y);
//       QUERY.eb(x, y, -1, -1, -1);
//     }
//   }

//   Dual_FenwickTree_2D<Monoid_Add<ll>, int, false> bit(X, Y);

//	for (auto& [a, b, c, d, x]: dat) bit.apply(a, c, b, d, x);
//	for (auto& [a, b, c, d, x]: QUERY) {
//	  if (x == -1) {
//	    ll ans = bit.get(a, b);
//	    print(ans);
//	  } else {
//	    bit.apply(a, c, b, d, x);
//	  }
//	}
//