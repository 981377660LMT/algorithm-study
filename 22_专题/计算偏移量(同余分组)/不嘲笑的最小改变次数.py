# !改变数组nums1的某些项 使得nums1的差分数组与nums2的差分数组一致
# 求最少移动次数
# !最后函数要平移到一起 考察两个数组的diff(力扣也有类似题)
# !diff的最大频率就是不用移动的元素个数

# 可以想像最后两个数组的函数图像，肯定是可以通过上下平移一个 delta 完全重合的。
# 现在要让不移动的点最多，就是计算平移之前各个点对应的纵坐标的差值delta，哪个数出现得最多。
# 最后这些点不需要移动，其他所有点都需要移动。

# 增量相同或者减量相同才不会互相嘲笑
# 修改TOM分数 至少多少关
from collections import Counter
import sys
from typing import List


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


def solve(nums1: List[int], nums2: List[int]) -> int:
    n = len(nums1)
    diff = [a - b for a, b in zip(nums1, nums2)]
    counter = Counter(diff)
    max_ = max(counter.values())
    return n - max_


n = int(input())
nums1 = list(map(int, input().split()))
nums2 = list(map(int, input().split()))


print(solve(nums1, nums2))
