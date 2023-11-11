# 求使得A变成B的操作的最小代价
# 操作1:加1/减1 代价为x
# 操作2:交换A中的相邻元素 代价为y
# 2<=n<=18

# 最暴力的方式:全排列然后计算对应的代价
# !希望固定每个Ai交换后对应的Bj
# 状压dp dp[index][state]表示nums2的前index个数 来自nums1中state`集合`时的最小代价 O(n*2^n)

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, cost1, cost2 = map(int, input().split())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))

    @lru_cache(None)
    def dfs(index: int, state: int) -> int:
        """
        dp[index][state]表示确定了index个数 哪几个位置的数对应交换后nums2 0,1,2...位置时的最小代价
        例如 0b1011 表示 a1,a2,a4 放到b1,b2,b3上的最小花费
        """
        if index == n:
            return 0

        res = INF
        inv = 0  # 邻位交换次数:产生的逆序对个数
        for cur in range(n - 1, -1, -1):
            if state & (1 << cur):
                inv += 1
                continue
            res = min(
                res,
                dfs(index + 1, state | (1 << cur))
                + abs(nums1[cur] - nums2[index]) * cost1
                + inv * cost2,  # 逆序对的个数
            )
        return res

    print(dfs(0, 0))

# より基本的な例題：
# ABC199E『Permutation』、 ABC180E『Traveling Salesman among Aerial Cities』
