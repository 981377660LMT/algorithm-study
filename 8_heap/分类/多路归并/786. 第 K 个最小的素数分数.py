from typing import List
from heapq import heappop, heappush

# 给你一个按递增顺序排序的数组 arr 和一个整数 k 。
# 数组 arr 由 1 和若干 素数  组成，且其中所有整数互不相同。
# 得到分数 arr[i] / arr[j] 。
# 第 k 个最小的分数是多少呢?  以长度为 2 的整数数组返回你的答案,
# 2 <= arr.length <= 1000
# 719. 找出第 k 小的距离对.py


class Solution:

    # 多路归并(val,row,col) 进出k次即可
    def kthSmallestPrimeFraction(self, arr: List[int], k: int) -> List[int]:
        n = len(arr)

        pq = []
        for j in range(1, n):
            heappush(pq, (arr[0] / arr[j], 0, j))

        for _ in range(k - 1):
            _, i, j = heappop(pq)
            if i + 1 < j:
                heappush(pq, (arr[i + 1] / arr[j], i + 1, j))

        _, i, j = pq[0]
        return [arr[i], arr[j]]


print(Solution().kthSmallestPrimeFraction(arr=[1, 2, 3, 5], k=3))
# 输出：[2,5]
# 解释：已构造好的分数,排序后如下所示:
# 1/5, 1/3, 2/5, 1/2, 3/5, 2/3
# 很明显第三个最小的分数是 2/5
