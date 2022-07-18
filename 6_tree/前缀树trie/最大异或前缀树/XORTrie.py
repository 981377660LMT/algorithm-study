from typing import List


class XORTrieNode:
    __slots__ = "bit", "count", "children"

    def __init__(self, bit: int):
        self.bit = bit  # 左右子树 bit=0表示左子树,bit=1表示右子树
        self.count = 0
        self.children: List["XORTrieNode"] = [None, None]  # type: ignore


class XORTrie:
    def __init__(self, bit=31):
        self.bit = bit
        self.root = XORTrieNode(-1)

    def insert(self, num: int) -> None:
        root = self.root
        for i in range(self.bit, -1, -1):
            bit = (num >> i) & 1
            if root.children[bit] is None:
                root.children[bit] = XORTrieNode(bit)
            root.children[bit].count += 1
            root = root.children[bit]

    def search(self, num: int) -> int:  # 查询，能获得的最大的异或值
        root = self.root
        res = 0
        for i in range(self.bit, -1, -1):
            bit = (num >> i) & 1
            needBit = bit ^ 1
            if root.children[needBit] is not None and root.children[needBit].count > 0:
                res = res << 1 | 1
                root = root.children[needBit]
            elif root.children[bit] is not None and root.children[bit].count > 0:
                res = res << 1
                root = root.children[bit]
        return res

    def discard(self, num: int) -> None:
        root = self.root
        for i in range(self.bit, -1, -1):
            if root is None:
                return
            bit = (num >> i) & 1
            root.children[bit].count -= 1
            root = root.children[bit]
