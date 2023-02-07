- `dp[j] = min{dp[i] + b[i]*a[j]}`
  如果 b[j]>=b[j+1](单调),那么可以用 CHT 优化 dp, `O(n^2) => O(n)`

- `dp[k][j] = min{dp[k-1][i] + b[k]*a[j]}`
  如果 b[j]>=b[j+1](单调),那么可以用 CHT 优化 dp, `O(n^2*k) => O(n*k)`

## 实现 CHT 的几种方法:

1. deque
   https://tjkendev.github.io/procon-library/python/convex_hull_trick/deque.html
   使用条件:
   - 追加的直线斜率单增/单减
   - 查询的最小值/最大值的 x 坐标单增/单减
2. 李超线段树
   https://tjkendev.github.io/procon-library/python/convex_hull_trick/li_chao_tree.html
   无法使用 deque 的时候,可以使用李超线段树，动态开点不需要预先知道查询，非常方便
