# from types import SimpleNamespace
from collections import namedtuple


def useXORTrie(bitLength=31):
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


if __name__ == '__main__':
    n = int(input())
    nums = list(map(int, input().split()))
    xorTire = useXORTrie()
    res = 0
    for num in nums:
        res = max(res, xorTire.search(num))
        xorTire.insert(num)
    print(res)

# https://stackoverflow.com/questions/1528932/how-to-create-inline-objects-with-properties
# python中像js一样创建对象
# 1. type
# res: IXORTrie = type('', (), {'insert': insert, 'search': search, 'discard': discard})
# 2. SimpleNamespace
# res: IXORTrie = SimpleNamespace(insert=insert, search=search, discard=discard)
# 3.namedtuple
# namedtuple('Res', ['insert', 'search', 'discard'])(insert, search, discard)
