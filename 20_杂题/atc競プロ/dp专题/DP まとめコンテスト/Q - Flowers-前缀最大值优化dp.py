"""移除花使得高度递增 求剩下的花的分数和最大值"""
# !即求分数最大的LIS 维护前缀的最大值
# n<=1e5
# 1<=h[i]<=n 且h[i]不重复


# 如果按照dp[index][height]的方式来做 看每个选不选 时间空间都是O(n^2)
# 需要改变dp定义
# !dp[i]表示高度为i结尾时的最大分数
# !dp[i]=max(dp[j])+h[i] (j<i)
# 因为只需要维护`前缀的最大值` 所以使用BIT


from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class BIT3:
    """单点修改 维护`前缀区间`最大值

    这么做正确的前提是不会删除或修改已经存进去的值
    每次都是加入新的值，这样已经存在的最大值一直有效。
    """

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
