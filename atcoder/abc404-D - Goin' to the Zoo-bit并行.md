# 题意简述

有 $N$ 个动物园，每个动物园有入场费 $C_i$。有 $M$ 种动物，每种动物可以在若干动物园看到。你可以多次访问同一个动物园。你要让每种动物都至少看两次，问最小总花费。

# 状态设计

- 你可以选择每个动物园访问 0、1、2 次（访问 2 次以上没有意义，因为每种动物最多只需要看 2 次）。
- 你需要记录每种动物已经看了几次（最多 2 次，超过 2 次也只算 2 次）。
- 由于 $M$ 可能较大，不能用 $M$ 维数组直接做 DP。

# 状态压缩

- 对每种动物 $i$，用 2bit 表示其观看次数 $G_i = \min(F_i, 2)$，$F_i$为实际观看次数。
- 所有动物的观看状态可以用 $2M$ bit 的整数 $watched$ 表示。
- $watched$ 的每 2bit 表示一种动物的观看状态（00=0 次，01=1 次，10=2 次）。

# 操作一：更新观看状态

- 进入某动物园时，该园能看到的动物 $i$，$G_i$ 需要加 1（但最大为 2）。
- 设 $one$ 是该动物园能看到的动物的掩码（每个动物 2bit，能看到的为 01，否则为 00）。
- 更新公式：$G_i \leftarrow \min(G_i + 1, 2)$
- 代码实现：`watched + (one & ~watched >> 1)`
  - $~watched >> 1$ 得到每个动物的“是否已经到 2 次”的掩码
  - $one & ~watched >> 1$ 只对还没到 2 次的动物加 1
  - 加法后，最多到 10（二进制 2）

# 操作二：判定是否所有动物都看了 2 次

- 只需判断 $watched$ 是否等于 $two\_all$，其中 $two\_all = \sum_{i=0}^{M-1} (2 << 2i)$，即每个动物的 2bit 都是 10（二进制 2）。

# DFS 搜索

- 对每个动物园，枚举访问 0、1、2 次的情况，递归搜索。
- 终止条件：所有动物都看了 2 次，返回当前花费，否则返回无穷大。
- 取三种选择的最小值。

# 代码关键片段讲解

```python
ones = [sum(1 << 2 * j for j in z) for z in zoo]
two_all = sum(2 << 2 * j for j in range(m))

def add_one(watched, one):
    return watched + (one & ~watched >> 1)

def dfs(i, watched, ans):
    if i == n:
        return ans if watched == two_all else 10**18
    ans0 = dfs(i + 1, watched, ans)  # 不选
    watched1 = add_one(watched, ones[i])
    ans1 = dfs(i + 1, watched1, ans + cost[i])  # 选一次
    watched2 = add_one(watched1, ones[i])
    ans2 = dfs(i + 1, watched2, ans + cost[i] * 2)  # 选两次
    return min(ans0, ans1, ans2)
```

- `ones[i]`：动物园 $i$ 能看到的动物的掩码
- `add_one(watched, ones[i])`：进入动物园 $i$ 后，更新观看状态
- `dfs(i+1, ...)`：递归到下一个动物园
- 终止条件：所有动物都看了 2 次（$watched == two\_all$）

# 为什么 add_one 正确？

- $one$ 是能看到的动物的掩码（每个动物 2bit，能看到的为 01）
- $watched$ 是当前状态（每个动物 2bit，00/01/10）
- $~watched >> 1$ 得到每个动物是否还没到 2 次（01 表示没到 2 次，00 表示已经 2 次）
- $one & ~watched >> 1$ 只对还没到 2 次的动物加 1
- $watched + (one & ~watched >> 1)$ 就是每个动物的观看次数最多加到 2

# 为什么判定所有动物都看了 2 次只需等于 two_all？

- $two\_all$ 是所有动物的 2bit 都为 10（二进制 2），即都看了 2 次
- 只需 $watched == two\_all$ 即可

# 总结

- 用 2bit 压缩每种动物的观看状态，整体用一个整数表示
- 每次递归枚举动物园访问次数，更新状态
- 用位运算高效更新和判定
- 总复杂度 $O(3^N \cdot M/W)$，$W$为机器字长，适合 $N \leq 15$，$M$较大时也能处理

如需进一步代码优化或并行化，可用记忆化搜索或多线程等手段。
