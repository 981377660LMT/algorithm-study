"""
对于查询 i ，请你找到 vali 和 pi 的 最大基因差 
其中 pi 是节点 nodei 到根之间的任意节点（包含 nodei 和根节点）
"""

from math import floor, log2
from typing import List
from collections import namedtuple
from typing import Protocol


class IXorTrie(Protocol):
    def insert(self, num: int) -> None:
        """将 `num` 插入到前缀树中"""
        ...

    def search(self, num: int) -> int:
        """查询 `num` 与前缀树中的最大异或值"""
        ...

    def discard(self, num: int) -> None:
        """在前缀树中删除 `num` 必须保证 `num` 在前缀树中存在"""
        ...


def useXORTrie(bitLength=31) -> "IXorTrie":
    trieRoot = [None, None, 0]

    def insert(num: int) -> None:
        root = trieRoot
        for i in range(bitLength, -1, -1):
            bit = (num >> i) & 1
            if root[bit] is None:  # type: ignore
                root[bit] = [None, None, 0]  # type: ignore
            root[bit][2] += 1  # type: ignore
            root = root[bit]  # type: ignore

    def search(num: int) -> int:
        root = trieRoot
        res = 0
        for i in range(bitLength, -1, -1):
            if root is None:  # Trie中未插入
                break

            bit = (num >> i) & 1
            needBit = bit ^ 1
            if root[needBit] is not None and root[needBit][2] > 0:  # type: ignore
                res = res << 1 | 1
                root = root[needBit]  # type: ignore
            else:
                res = res << 1
                root = root[bit]  # type: ignore

        return res

    def discard(num: int) -> None:
        root = trieRoot
        for i in range(bitLength, -1, -1):
            if root is None:  # Trie中未插入
                break

            bit = (num >> i) & 1
            if root[bit] is not None:  # type: ignore
                root[bit][2] -= 1  # type: ignore
            root = root[bit]  # type: ignore

    return namedtuple("XORTrie", ["insert", "search", "discard"])(insert, search, discard)  # type: ignore


class Solution:
    def maxGeneticDifference(self, parents: List[int], queries: List[List[int]]) -> List[int]:
        def dfs(cur: int) -> None:
            trie.insert(cur)
            for qi, qv in nodeQueries[cur]:
                res[qi] = trie.search(qv)
            for next in adjList[cur]:
                dfs(next)
            trie.discard(cur)

        n = len(parents)
        root = -1
        adjList = [[] for _ in range(n)]
        for i in range(n):
            if parents[i] != -1:
                adjList[parents[i]].append(i)
            else:
                root = i

        nodeQueries = [[] for _ in range(n)]
        for i, (node, val) in enumerate(queries):
            nodeQueries[node].append((i, val))

        bit = floor(log2(int(2e5))) + 1  # !vi<=2e5
        trie = useXORTrie(bit)

        res = [-1] * len(queries)
        dfs(root)
        return res


if __name__ == "__main__":
    print(
        Solution().maxGeneticDifference(
            parents=[3, 7, -1, 2, 0, 7, 0, 2], queries=[[4, 6], [1, 15], [0, 5]]
        )
    )
