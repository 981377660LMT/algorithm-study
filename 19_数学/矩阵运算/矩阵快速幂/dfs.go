// int id[maxn];
// void dfs(int q, int k) {
//     if (q == n) return calc();
//     id[q] = k + 1, dfs(q + 1, k + 1);
//     for (int i = 1; i <= k; ++i) id[q] = i, dfs(q + 1, k);
// }