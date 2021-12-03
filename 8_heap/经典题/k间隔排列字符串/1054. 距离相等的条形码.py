from collections import Counter
from typing import List

# 两个相邻的条形码 不能 相等
# 此题保证存在答案。
class Solution:
    def rearrangeBarcodes(self, barcodes: List[int]) -> List[int]:
        counter = Counter(barcodes)
        barcodes.sort(key=lambda num: (counter[num], num))
        barcodes[1::2], barcodes[::2] = (
            barcodes[: len(barcodes) >> 1],
            barcodes[len(barcodes) >> 1 :],
        )
        return barcodes


print(Solution().rearrangeBarcodes([1, 1, 1, 2, 2, 2]))
print(Solution().rearrangeBarcodes([4, 3, 8, 4, 4, 4, 8, 3, 3, 3]))
