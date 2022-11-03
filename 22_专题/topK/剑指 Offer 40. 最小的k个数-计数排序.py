# 0 <= k <= arr.length <= 10000
# 0 <= arr[i] <= 10000
# !最小的k个数 计数排序
from typing import List


class Solution:
    def getLeastNumbers(self, arr: List[int], k: int) -> List[int]:
        if k == 0:
            return []

        counter = [0] * 10010
        for num in arr:
            counter[num] += 1

        res = []
        for i in range(10010):
            if counter[i] > 0:
                res.extend([i] * counter[i])
                if len(res) >= k:
                    break

        return res[:k]
