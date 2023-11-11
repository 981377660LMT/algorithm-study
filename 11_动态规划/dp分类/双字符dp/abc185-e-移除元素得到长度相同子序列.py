# https://atcoder.jp/contests/abc185/tasks/abc185_e

# 输入 n(≤1000) 和 m(≤1000)，长度分别为 n 和 m 的数组 a 和 b，元素范围 [1,1e9]。
# 从 a 中移除元素，得到一个子序列 a'；从 b 中移除若干元素，得到一个子序列 b'。
# 要求 a' 和 b' 的长度相同。
# 输出 (a和b总共移除的元素个数) + (a'[i]≠b'[i]的i的个数) 的最小值。

# https://atcoder.jp/contests/abc185/submissions/36352936

# 类似 LCS，定义 f[i][j] 表示 a 的前 i 个元素和 b 的前 j 的元素算出的答案。

# - 不选 a[i] 选 b[j]：f[i-1][j]+1
# - 选 a[i] 不选 b[j]：f[i][j-1]+1
# - 选 a[i] 选 b[j]：f[i-1][j-1] + (a[i]==b[j] ? 0 : 1)
# 取最小值。

# 注：都不选是不用考虑的，这已经包含在 f[i-1][j] 或者 f[i][j-1] 中了。
# 也可以这么想：都不选是不如都选的。

# 边界：f[0][i]=f[i][0]=i。
# 答案：f[n][m]。

from functools import lru_cache
from typing import List


def sequnceMatching1(nums1: List[int], nums2: List[int]) -> int:
    n, m = len(nums1), len(nums2)
    dp = [[0] * (m + 1) for _ in range(n + 1)]
    for i in range(n + 1):
        dp[i][0] = i
    for j in range(m + 1):
        dp[0][j] = j
    for i in range(1, n + 1):
        for j in range(1, m + 1):
            dp[i][j] = min(
                dp[i - 1][j] + 1,
                dp[i][j - 1] + 1,
                dp[i - 1][j - 1] + (nums1[i - 1] != nums2[j - 1]),
            )
    return dp[n][m]


def sequnceMatching2(nums1: List[int], nums2: List[int]) -> int:
    @lru_cache(None)
    def dfs(i: int, j: int) -> int:
        if i == n:
            return m - j
        if j == m:
            return n - i
        return min(
            dfs(i + 1, j) + 1,
            dfs(i, j + 1) + 1,
            dfs(i + 1, j + 1) + (nums1[i] != nums2[j]),
        )

    n, m = len(nums1), len(nums2)
    return dfs(0, 0)


if __name__ == "__main__":

    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353
    INF = int(4e18)

    n, m = map(int, input().split())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    print(sequnceMatching1(nums1, nums2))
    # print(sequnceMatching2(nums1, nums2))
