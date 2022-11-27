from typing import List


class XORTrieNode:
    __slots__ = "bit", "count", "children", "index"

    def __init__(self):
        self.count = 0
        self.children: List["XORTrieNode"] = [None, None]  # type: ignore


class XORTrie:
    def __init__(self, upper: int):
        self.bit = upper.bit_length()
        self.root = XORTrieNode()

    def insert(self, num: int) -> None:
        root = self.root
        for i in range(self.bit, -1, -1):
            bit = (num >> i) & 1
            if root.children[bit] is None:
                root.children[bit] = XORTrieNode()
            root.children[bit].count += 1
            root = root.children[bit]

    def search(self, num: int, upper: int) -> int:
        """
        求num与树中异或小于等于upper的个数
        """
        root = self.root
        res = 0
        for i in range(self.bit, -1, -1):
            if root is None:
                return res
            bit = (num >> i) & 1
            bitLimit = (upper >> i) & 1
            if bitLimit == 1:
                if root.children[bit] is not None:
                    res += root.children[bit].count
                root = root.children[bit ^ 1]
            else:
                root = root.children[bit]
        return res


class Solution:
    def countPairs(self, nums: List[int], low: int, high: int) -> int:
        trie = XORTrie(int(2e4))
        res = 0
        for num in nums:
            res += trie.search(num, high + 1) - trie.search(num, low)  # !注意这里是high+1
            trie.insert(num)
        return res
