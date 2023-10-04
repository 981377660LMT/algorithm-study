# 设计一个数据结构,有效地找到给定子数组的 出现 threshold 次数或次数以上的元素 。
# 1 <= arr.length <= 2e4
# threshold <= right - left + 1
# 2 * threshold > right - left + 1
# 调用 query 的次数最多为 1e4

# !根号分治
# !针对不同的询问区间长度，使用两种不同的方法。
# 记 SQRT = sqrt(2n)
# - 区间长度小于 SQRT ，使用暴力计算
# - 区间长度大于 SQRT ，则绝对众数出现次数 大于 SQRT/2
#   可能的候选人个数不超过 2n/SQRT ，使用前缀和统计频率大于SQRT/2的数的出现次数


from collections import Counter, defaultdict
from typing import List


class MajorityChecker:
    def __init__(self, arr: List[int]):
        n = len(arr)
        self.SQRT = int(2 * n**0.5)
        self.arr = arr
        self.preSumMap = defaultdict(lambda: [0])  # !只统计频率大于SQRT/2的数的前缀和
        counter = Counter(arr)
        for num, count in counter.items():
            if count > self.SQRT // 2:
                preSum = [0] * (n + 1)
                for i in range(n):
                    preSum[i + 1] = preSum[i] + (arr[i] == num)
                self.preSumMap[num] = preSum

    def query(self, left: int, right: int, threshold: int) -> int:
        """寻找闭区间[left,right]中的绝对众数,调用 query 的次数最多为 1e4"""
        length = right - left + 1
        if length <= self.SQRT:
            counter = dict()
            for i in range(left, right + 1):
                num = self.arr[i]
                counter[num] = counter.get(num, 0) + 1
                if counter[num] >= threshold:
                    return num
            return -1
        else:
            for num, preSum in self.preSumMap.items():
                if preSum[right + 1] - preSum[left] >= threshold:
                    return num
            return -1
