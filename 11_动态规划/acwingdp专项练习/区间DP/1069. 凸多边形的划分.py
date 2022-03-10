# 给定一个具有 N 个顶点的凸多边形，将顶点从 1 至 N 标号，每个顶点的权值都是一个正整数。
# 将这个凸多边形划分成 N−2 个互不相交的三角形，对于每个三角形，其三个顶点的权值相乘都可得到一个权值乘积，试求所有三角形的顶点权值乘积之和至少为多少。
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))

n = int(input())
nums = list(map(int, input().split()))


@lru_cache(maxsize=None)
def dfs(left: int, right: int) -> int:
    """[left:right+1]这一段合并的代价最小"""
    """i-j连起来的这条边`不动`，枚举根这条边配对的点的位置，枚举切割位置"""
    if right - left < 2:
        return 0
    if right - left == 2:
        return nums[left] * nums[left + 1] * nums[right]

    res = int(1e100)
    for i in range(left + 1, right):
        res = min(res, dfs(left, i) + dfs(i, right) + nums[left] * nums[i] * nums[right])
    return res


print(dfs(0, n - 1))
