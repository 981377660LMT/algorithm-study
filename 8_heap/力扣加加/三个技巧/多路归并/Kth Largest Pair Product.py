# 两个数组第k个最大的对的乘积
from heapq import merge
from itertools import islice
from typing import Generator, List


class Solution:
    def solve(self, a: List[int], b: List[int], k: int):
        def g1(index: int) -> Generator[int, None, None]:
            return (a[index] * num for num in b)

        def g2(index: int) -> Generator[int, None, None]:
            return (a[index] * num for num in reversed(b))

        a, b = sorted(a), sorted(b)
        # 生成器的多路
        gen = (g2(i) if num > 0 else g1(i) for i, num in enumerate(a))
        iter = merge(*gen, reverse=True)
        return next(islice(iter, k, None))  # 切片，第k个最大


print(Solution().solve(a=[-2, 4, 3], b=[5, 7], k=2))
# The three largest products are 4 * 7 = 28, 3 * 7 = 21 and 4 * 5 = 20.
