from collections import defaultdict
from typing import List

# 1 <= A.length <= 1000
# nums[i]<=2^16
# 找出索引为 (i, j, k) 的三元组
# A[i] & A[j] & A[k] == 0，其中 & 表示按位与（AND）操作符。
class Solution:
    def countTriplets(self, A: List[int]) -> int:
        memo = defaultdict(int)
        for n1 in A:
            for n2 in A:
                memo[n1 & n2] += 1

        res = 0
        for num in A:
            for key, val in memo.items():
                if num & key == 0:
                    res += val
        return res

    def countTriplets2(self, A: List[int]) -> int:
        counter = defaultdict(int)
        for n1 in A:
            for n2 in A:
                counter[n1 & n2] += 1

        res = 0
        # 遍历每个数字去找满足条件的数量
        for num in A:
            for key, val in counter.items():
                if num & key == 0:
                    res += val
        return res
