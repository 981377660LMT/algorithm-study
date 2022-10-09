# 设计一个数据结构,有效地找到给定子数组的 出现 threshold 次数或次数以上的元素 。
# 1 <= arr.length <= 2e4
# threshold <= right - left + 1
# 2 * threshold > right - left + 1
# 调用 query 的次数最多为 1e4

from bisect import bisect_left, bisect_right
from collections import defaultdict
from random import randint
from typing import List


class MajorityChecker:
    def __init__(self, arr: List[int]):
        self.arr = arr
        self.indexes = defaultdict(list)  # 预处理,便于O(logn)查询区间内数字的频率
        for i, num in enumerate(arr):
            self.indexes[num].append(i)

    def query(self, left: int, right: int, threshold: int) -> int:
        """寻找闭区区间[left,right]中的绝对众数,调用 query 的次数最多为 1e4

        1. 如果一个区间内确实存在绝对众数,那么我们随机选择一个数,
        这个数为绝对众数的概率至少是 1/2 。
        2. 我们可以计算出这个数在区间内出现了多少次 (可以预处理+二分 O(logn)求出) 。
        如果出现频率大于等于 threshold ,那么这个数就是区间的绝对众数。
        3. 如果随机选择了20次,都没有找到符合条件的数,那么就说明区间内确实不存在绝对众数。

        误判的概率为1/2^20,即确实存在绝对众数,但是抽取了20次都没有抽到绝对众数

        # !时间复杂度为 O(q*logn*20) ,其中 q 为 query 的次数。
        """

        ROUND = 20
        for _ in range(ROUND):
            index = randint(left, right)
            num = self.arr[index]
            count = bisect_right(self.indexes[num], right) - bisect_left(self.indexes[num], left)
            if count >= threshold:
                return num
        return -1
