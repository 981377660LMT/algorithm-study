from typing import List


class Solution:
    def checkIfExist(self, arr: List[int]) -> bool:
        record = {v * 2: i for i, v in enumerate(arr)}
        for i, v in enumerate(arr):
            if v in record and i != record[v]:
                return True
        return False


# 输入：arr = [7,1,14,11]
# 输出：true
# 解释：N = 14 是 M = 7 的两倍，即 14 = 2 * 7 。
