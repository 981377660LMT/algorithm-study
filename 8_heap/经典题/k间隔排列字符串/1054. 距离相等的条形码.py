from collections import Counter
from typing import List

# 两个相邻的条形码 不能 相等
# 此题保证存在答案。

# 按counter.most_common()将字符push进，然后half切片前一半大
class Solution:
    def rearrangeBarcodes(self, barcodes: List[int]) -> List[int]:
        counter = Counter(barcodes)
        chars = []
        for char, count in counter.most_common():
            for _ in range(count):
                chars.append(char)
        half = (len(chars) + 1) >> 1

        left, right = chars[:half], chars[half:]

        res = [0] * len(chars)
        res[::2], res[1::2] = left, right

        return res


print(Solution().rearrangeBarcodes([1, 1, 1, 2, 2, 2]))
print(Solution().rearrangeBarcodes([4, 3, 8, 4, 4, 4, 8, 3, 3, 3]))
print(Solution().rearrangeBarcodes([2, 1, 1]))
