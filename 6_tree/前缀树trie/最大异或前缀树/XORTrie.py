from typing import List


class XORTrie:
    def __init__(self, bitLength=31):
        self.children: List[XORTrie] = [None, None]  # type: ignore
        self.bitLength = bitLength

    def insert(self, x: int) -> None:
        root = self
        for i in range(self.bitLength, -1, -1):
            bit = (x >> i) & 1
            if root.children[bit] is None:
                root.children[bit] = XORTrie()
            root = root.children[bit]

    def query(self, x: int) -> int:  # 查询，能获得的最大的异或值
        root = self
        res = 0
        for i in range(self.bitLength, -1, -1):
            bit = (x >> i) & 1
            needBit = bit ^ 1
            if root.children[needBit] is not None:
                res = res << 1 | 1
                root = root.children[needBit]
            else:
                res = res << 1
                root = root.children[bit]
        return res

