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


def useXORTrie(upper: int) -> "IXorTrie":
    trieRoot = [None, None, 0]
    bitLength = upper.bit_length()

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
                res |= 1 << i
                root = root[needBit]  # type: ignore
            else:
                root = root[bit]  # type: ignore

        return res

    def remove(num: int) -> None:
        root = trieRoot
        for i in range(bitLength, -1, -1):
            if root is None:
                raise ValueError(f"fail to remove: num {num} not in trie")

            bit = (num >> i) & 1
            if root[bit] is not None:  # type: ignore
                root[bit][2] -= 1  # type: ignore
            root = root[bit]  # type: ignore

    return namedtuple("XORTrie", ["insert", "search", "discard"])(insert, search, remove)  # type: ignore


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    xorTire = useXORTrie(int(1e9 + 10))
    res = 0
    for num in nums:
        res = max(res, xorTire.search(num))
        xorTire.insert(num)
    print(res)

# https://stackoverflow.com/questions/1528932/how-to-create-inline-objects-with-properties
# python中像js一样创建对象
# !1. type 但是很慢
# res: IXORTrie = type('', (), {'insert': insert, 'search': search, 'discard': discard})
# !2. SimpleNamespace
# res: IXORTrie = SimpleNamespace(insert=insert, search=search, discard=discard)
# !3.namedtuple
# namedtuple('Res', ['insert', 'search', 'discard'])(insert, search, discard)

# !simplenamespace vs namedtuple
