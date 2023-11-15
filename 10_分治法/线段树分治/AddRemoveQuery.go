// https://maspypy.github.io/library/ds/offline_query/add_remove_query.hpp

package main

func main() {

}

// /*
// ・時刻 t に x を追加する
// ・時刻 t に x を削除する
// があるときに、
// ・時刻 [l, r) に x を追加する
// に変換する。
// クエリが時系列順に来ることが分かっているときは monotone = true の方が高速。
// */
// template <typename X, bool monotone>
// struct Add_Remove_Query {
//   map<X, int> MP;
//   vc<tuple<int, int, X>> dat;
//   map<X, vc<int>> ADD;
//   map<X, vc<int>> RM;

//   void add(int time, X x) {
//     if (monotone) return add_monotone(time, x);
//     ADD[x].eb(time);
//   }
//   void remove(int time, X x) {
//     if (monotone) return remove_monotone(time, x);
//     RM[x].eb(time);
//   }

//   // すべてのクエリが終わった現在時刻を渡す
//   vc<tuple<int, int, X>> calc(int time) {
//     if (monotone) return calc_monotone(time);
//     vc<tuple<int, int, X>> dat;
//     for (auto&& [x, A]: ADD) {
//       vc<int> B;
//       if (RM.count(x)) {
//         B = RM[x];
//         RM.erase(x);
//       }
//       if (len(B) < len(A)) B.eb(time);
//       assert(len(A) == len(B));

//       sort(all(A));
//       sort(all(B));
//       FOR(i, len(A)) {
//         assert(A[i] <= B[i]);
//         if (A[i] < B[i]) dat.eb(A[i], B[i], x);
//       }
//     }
//     assert(len(RM) == 0);
//     return dat;
//   }

// private:
//   void add_monotone(int time, X x) {
//     assert(!MP.count(x));
//     MP[x] = time;
//   }
//   void remove_monotone(int time, X x) {
//     auto it = MP.find(x);
//     assert(it != MP.end());
//     int t = (*it).se;
//     MP.erase(it);
//     if (t == time) return;
//     dat.eb(t, time, x);
//   }
//   vc<tuple<int, int, X>> calc_monotone(int time) {
//     for (auto&& [x, t]: MP) {
//       if (t == time) continue;
//       dat.eb(t, time, x);
//     }
//     return dat;
//   }
// };
