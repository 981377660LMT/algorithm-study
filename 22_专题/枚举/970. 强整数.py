# 给定两个正整数 x 和 y，如果某一整数等于 x^i + y^j，
# 其中整数 i >= 0 且 j >= 0，那么我们认为该整数是一个强整数。

# 返回值小于或等于 bound 的所有强整数组成的列表。
from typing import List

# 1 <= x <= 100
# 1 <= y <= 100
# 0 <= bound <= 10^6
class Solution:
    def powerfulIntegers(self, x: int, y: int, bound: int) -> List[int]:
        xs = {x ** i for i in range(20) if x ** i < bound}
        ys = {y ** i for i in range(20) if y ** i < bound}
        return list({i + j for i in xs for j in ys if i + j <= bound})


print(Solution().powerfulIntegers(x=2, y=3, bound=10))
# 输出：[2,3,4,5,7,9,10]
# 解释：
# 2 = 2^0 + 3^0
# 3 = 2^1 + 3^0
# 4 = 2^0 + 3^1
# 5 = 2^1 + 3^1
# 7 = 2^2 + 3^1
# 9 = 2^3 + 3^0
# 10 = 2^0 + 3^2
