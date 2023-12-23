# https://leetcode.cn/problems/count-paths-that-can-form-a-palindrome-in-a-tree/solution/dfsyu-chu-li-gen-jie-dian-dao-mei-ge-jie-focq/
from typing import Counter, List


class Solution:
    def countPalindromePaths(self, parent: List[int], s: str) -> int:
        n = len(s)
        adjList = [[] for _ in range(n)]
        for i in range(1, n):
            p = parent[i]
            mask = 1 << (ord(s[i]) - 97)
            adjList[i].append((p, mask))
            adjList[p].append((i, mask))

        xorToRoot = [0] * n  # 节点 i 到根节点的异或和

        def dfs(cur: int, pre: int) -> None:
            for next, weight in adjList[cur]:
                if next == pre:
                    continue
                xorToRoot[next] = xorToRoot[cur] ^ weight  # !op
                dfs(next, cur)

        dfs(0, -1)

        res = 0
        counter = Counter(xorToRoot)
        for mask, count in counter.items():
            res += count * (count - 1)
            for c in range(26):
                res += counter[mask ^ (1 << c)] * count
        return res // 2
