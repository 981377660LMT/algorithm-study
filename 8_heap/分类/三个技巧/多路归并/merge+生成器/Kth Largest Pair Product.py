# 两个数组第k个最大的对的乘积
from heapq import heappop, heappush, merge
from itertools import islice
from typing import Generator, List


class Solution:
    def solve(self, a: List[int], b: List[int], k: int):
        def gen1(index: int) -> Generator[int, None, None]:
            return (a[index] * num for num in b)

        def gen2(index: int) -> Generator[int, None, None]:
            return (a[index] * num for num in reversed(b))

        a, b = sorted(a), sorted(b)
        # 生成器的多路
        allGen = [gen2(i) if num > 0 else gen1(i) for i, num in enumerate(a)]
        iterable = merge(*allGen, reverse=True)
        return next(islice(iterable, k, None))  # 切片，第k个最大


print(Solution().solve(a=[-2, 4, 3], b=[5, 7], k=2))
print(Solution().solve(a=[-4, -2, 0, 3], b=[2, 4], k=5))
# The three largest products are 4 * 7 = 28, 3 * 7 = 21 and 4 * 5 = 20.


class Solution2:
    def solve(self, a, b, k):
        a = sorted(a)
        b = sorted(b)
        k += 1
        heap = []
        for val in b:
            if val >= 0:
                for i in range(len(a) - 1, -1, -1):
                    cand = val * a[i]
                    if heap and len(heap) == k and cand <= heap[0]:
                        break
                    heappush(heap, cand)
                    if len(heap) > k:
                        heappop(heap)
            else:
                for i in range(len(a)):
                    cand = val * a[i]
                    if heap and len(heap) == k and cand <= heap[0]:
                        break
                    heappush(heap, cand)
                    if len(heap) > k:
                        heappop(heap)
        return heap[0]
