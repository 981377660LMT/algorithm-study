from typing import List

# 请你返回一个最小整数 value ，使得将数组中所有大于 value 的值变成 value 后，
# 数组的和最接近  target （最接近表示两者之差的绝对值最小）。


class Solution:
    def findBestValue(self, arr: List[int], target: int) -> int:
        ...


print(Solution().findBestValue(arr=[4, 9, 3], target=10))
