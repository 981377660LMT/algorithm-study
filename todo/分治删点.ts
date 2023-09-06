// // 调用g时，s为对除了i以外所有点均调用过了f的状态。但不保证调用f的顺序
// // 总计会调用 $O(NlgN)$ 次的 f和g

// void f(State& s, int i);
// void g(const State& s, int i);

// template <typename State, typename F, typename G>
// void divideconquer(const State& s, int l, int r, const F& f, const G& g) {
//     if (r == l + 1) {
//         g(s, l);
//         return;
//     }
//     int m = (r + l) / 2;

//     {
//         State copy = s;
//         for (int i = l; i < m; i++) {
//             f(copy, i);
//         }
//         divideconquer(copy, m, r, f, g);
//     }

//     {
//         State copy = s;
//         for (int i = m; i < r; i++) {
//             f(copy, i);
//         }
//         divideconquer(copy, l, m, f, g);
//     }
// }

function divideConquer(): void {}
