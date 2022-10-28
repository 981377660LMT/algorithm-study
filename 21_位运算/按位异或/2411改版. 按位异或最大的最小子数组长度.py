# 对于 0 到 n - 1 之间的每一个下标 i ，
# 你需要找出 nums 中一个 最小 非空子数组，
# 它的起始位置为 i （包含这个位置），同时有 最大 的 按位异或运算值 。

# 倒序遍历，求后缀异或和
# 01Trie的结点多存储一个index信息
# 如果有多个相同的，存下标最小的


import random
from typing import List

INF = int(1e18)


class XORTrieNode:
    __slots__ = "children", "index"

    def __init__(self):
        self.children: List["XORTrieNode"] = [None, None]  # type: ignore
        self.index = -1


class XORTrie:
    def __init__(self, bit=31):
        self.bit = bit
        self.root = XORTrieNode()

    def insert(self, num: int, index: int) -> None:
        root = self.root
        for i in range(self.bit, -1, -1):
            bit = (num >> i) & 1
            if root.children[bit] is None:
                root.children[bit] = XORTrieNode()
            root = root.children[bit]
            if root.index == -1 or root.index > index:
                root.index = index

    def search(self, num: int) -> int:  # !查询能获得的最大的异或值时的最小下标
        root = self.root
        res = INF
        for i in range(self.bit, -1, -1):
            bit = (num >> i) & 1
            needBit = bit ^ 1
            if root.children[needBit] is not None:
                root = root.children[needBit]
            elif root.children[bit] is not None:
                root = root.children[bit]
            res = root.index
        return res


class Solution:
    def smallestSubarrays(self, nums: List[int]) -> List[int]:
        n = len(nums)
        res = [1] * n
        trie = XORTrie()
        sufXor = 0
        for i in range(n - 1, -1, -1):
            trie.insert(sufXor, i)
            sufXor ^= nums[i]
            res[i] = trie.search(sufXor) - i + 1
        return res

    def smallestSubarrays2(self, nums: List[int]) -> List[int]:
        """暴力对拍"""
        n = len(nums)
        res = [1] * n
        for i in range(n):
            xor_, curMax, curLen = 0, 0, 0
            for j in range(i, n):
                xor_ ^= nums[j]
                if xor_ > curMax:
                    curMax = xor_
                    curLen = j - i + 1
                    res[i] = curLen
        return res


# random test
for _ in range(10000):
    n = random.randint(1, 10)
    nums = [random.randint(0, int(1e9)) for _ in range(n)]
    assert Solution().smallestSubarrays(nums) == Solution().smallestSubarrays2(nums)
