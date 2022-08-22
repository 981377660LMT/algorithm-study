"""
https://leetcode.cn/problems/find-the-k-sum-of-an-array/solution/by-tsreaper-ps7w/
多路归并求数组的第k大的子序列和
k,n<=1e5
O(nlogn + klogk)

转化为多路归并母题
给定n个非负数a1, a2,...,an,求第k个最小的子序列和。

几个技巧:
1. 图中遍历子序列的方法 
可以理解为多路归并/dijkstra遍历图求最短路 
这里全部加上负数可以理解为dijk在处理负权边的时候需要一个势能函数,最后求出来的最短路减去这个势能

2. 引入负数
将负数转化为绝对值后
此时只需要在答案上减去所有负数的绝对值即可
因为此时选负数等价于不选变化后的正数
不选负数对应选变化后的正数
这种技巧可见 leetcode 1982 - 从子集的和还原数组 https://leetcode.cn/problems/find-array-given-subset-sums/。

3. 最大子序列和变为最小子序列和
因为第k大子序列和对应第k小子序列和 (原来取的变为全不取,取的变为全取 即 sumK => allSum - sumK)


至此 原问题变为
n个非负整数,求第k小的子序列和,k最大2000
"""

from collections import deque
from heapq import heappop, heappush
from typing import List


def kSum(nums: List[int], k: int) -> int:
    """求数组中第k个最大的子序列和"""
    n = len(nums)
    allSum, negSum = 0, 0
    for i in range(n):
        allSum += nums[i]
        if nums[i] < 0:
            negSum += nums[i]
            nums[i] = -nums[i]

    nums.sort()
    res = 0  # k = 1 时的答案，也就是空集的和
    pq = [(nums[0], 0)]  # 堆，元素是 (子序列和，最后一个元素的索引)
    for _ in range(k - 1):
        cur, index = heappop(pq)
        res = cur
        if index + 1 < n:
            heappush(pq, (cur + nums[index + 1], index + 1))
            heappush(pq, (cur - nums[index] + nums[index + 1], index + 1))

    res += negSum  # 消除负数影响
    return allSum - res  # 变为第k大的子序列和


print(kSum(nums=[1, -2, 3, 4, -10, 12], k=16))

# 遍历子序列
def traverseSequnce(s: str) -> List[str]:
    """
    遍历子序列的一种方式
    每个子序列视作一个结点
    from = s1s2...si
    则有向边为
    1. from -> s1s2...sisi+1  添加后继
    2. from -> s1s2...si-1si+1 弹出前驱，添加后继

    这个遍历方法相当于根据子序列最后一个元素把子序列分为了若干个类，
    并规定只有相邻的类之间有边，大大减少了图中边的数量，
    还确保了从源点到任何一个点只有唯一的一条路径。

    这种建图好处是可以沿着子序列某个性质上升或者下降的方向遍历图
    配合堆 使得每次放进去的数不小于/不大于拿出来的数
    """
    if not s:
        return [""]
    n, res = len(s), [""]
    queue = deque([(s[0], 0)])  # 子序列为 s, 最后一个元素是第 i 个元素的子序列。
    while queue:
        cur, index = queue.popleft()
        res.append(cur)
        if index + 1 < n:
            queue.append((cur + s[index + 1], index + 1))
            queue.append((cur[:-1] + s[index + 1], index + 1))
    return res


print(traverseSequnce("abc"))
