"""移除花使得高度递增 求剩下的花的分数和最大值"""
# !即求分数最大的LIS 维护前缀的最大值
# n<=1e5
# 1<=h[i]<=n 且h[i]不重复
# scores[i]<=1e9

# 如果按照dp[index][height]的方式来做 看每个选不选 时间空间都是O(n^2)
# 需要分析转移式优化
# !选第index个花的时候 唯一的影响因素为选了的前一个花的高度
# !dp[i][j] = max(dp[i-1][k] for k in raneg(j)) + h[i]   (h[i]==j 时)
# !dp[i][j] = dp[i-1][j]   (h[i]!=j 时)
# 因为只需要维护`前缀的最大值` 所以使用BIT


from collections import defaultdict
import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class BIT3:
    """单点修改 维护`前缀区间`最大值"""

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    def update(self, index: int, target: int) -> None:
        """将后缀区间`[index,size]`的最大值更新为target"""
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] = max(self.tree[index], target)
            index += index & -index

    def query(self, index: int) -> int:
        """查询前缀区间`[1,index]`的最大值"""
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res = max(res, self.tree[index])
            index -= index & -index
        return res


##########################################################################
n = int(input())
heights = list(map(int, input().split()))
scores = list(map(int, input().split()))

bit = BIT3(n)
for i in range(n):
    preMax = bit.query(heights[i])
    bit.update(heights[i], preMax + scores[i])

print(bit.query(n))
