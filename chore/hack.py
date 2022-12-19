# generate random cases, until find a case that solution1 is not equal to solution2


from typing import List


class Solution1:
    def isPossible(self, n: int, edges: List[List[int]]) -> bool:
        ...


class Solution2:
    def isPossible(self, n: int, edges: List[List[int]]) -> bool:
        ...


import random

while True:
    n = random.randint(1, 100)
    edges = []
    for _ in range(random.randint(0, n * (n - 1) // 2)):
        u, v = random.randint(1, n), random.randint(1, n)
        if u != v and [u, v] not in edges and [v, u] not in edges:
            edges.append([u, v])
    res1 = Solution1().isPossible(n, edges)
    res2 = Solution2().isPossible(n, edges)
    if res1 != res2:
        print(n, edges)
        print(res1, res2)
        break
