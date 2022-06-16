from collections import namedtuple
from typing import List


def useArrayXORTrie(bitLength=31):
    trieRoot = [None, None, 0]

    def insert(num: int) -> None:
        root = trieRoot
        for i in range(bitLength, -1, -1):
            bit = (num >> i) & 1
            if root[bit] is None:
                root[bit] = [None, None, 0]
            root[bit][2] += 1
            root = root[bit]

    def search(num: int) -> int:
        root = trieRoot
        res = 0
        for i in range(bitLength, -1, -1):
            if root is None:  # Trie中未插入
                break

            bit = (num >> i) & 1
            needBit = bit ^ 1
            if root[needBit] is not None and root[needBit][2] > 0:
                res = res << 1 | 1
                root = root[needBit]
            else:
                res = res << 1
                root = root[bit]

        return res

    def discard(num: int) -> None:
        root = trieRoot
        for i in range(bitLength, -1, -1):
            if root is None:  # Trie中未插入
                break

            bit = (num >> i) & 1
            if root[bit] is not None:
                root[bit][2] -= 1
            root = root[bit]

    return namedtuple('XORTrie', ['insert', 'search', 'discard'])(insert, search, discard)


class Solution:
    def maximizeXor(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        # offline，排序，Trie树 贪心
        n = len(nums)
        nums.sort()

        qn = len(queries)
        Q = [(qv, limit, qi) for qi, (qv, limit) in enumerate(queries)]
        Q.sort(key=lambda x: x[1])

        res = [0] * qn
        xorTrie = useArrayXORTrie()
        ni = 0

        for qv, limit, qi in Q:
            while ni < n and nums[ni] <= limit:  # 不超过m的都入树
                xorTrie.insert(nums[ni])
                ni += 1

            if ni == 0:  # 树中没有元素
                res[qi] = -1
            else:
                res[qi] = xorTrie.search(qv)

        return res
