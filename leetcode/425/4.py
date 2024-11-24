from typing import List


class Solution:
    def maximizeSumOfWeights(self, edges: List[List[int]], k: int) -> int:

        edges.sort(key=lambda x: -x[2])

        degrees = [0] * (max(max(u, v) for u, v, _ in edges) + 1)
        total = 0

        for u, v, w in edges:
            if degrees[u] < k and degrees[v] < k:
                total += w
                degrees[u] += 1
                degrees[v] += 1

        return total
