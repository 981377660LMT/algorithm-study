# 每个查询询问
# !nums1的前x项的集合是否与nums2的前y项的集合相等 (重复元素只算一次)
# n<=2e5
# numsi<=1e9
# x,y<=n


# !异或哈希/异或前缀和 用随机数产生哈希值 用异或来计算区间所含集合的哈希值
# Zobrist Hash
# !xorではなく和を使うと個数に対応したハッシュが作れる


from collections import defaultdict
from random import randint


class FastHashSet:
    """快速计算哈希值的集合."""

    _poolSingleton = defaultdict(lambda: randint(1, (1 << 61) - 1))

    __slots__ = ("_set", "_hash")

    def __init__(self) -> None:
        self._set = set()
        self._hash = 0

    def add(self, x: int) -> None:
        if x not in self._set:
            self._set.add(x)
            self._hash ^= self._poolSingleton[x]

    def discard(self, x: int) -> None:
        if x in self._set:
            self._set.discard(x)
            self._hash ^= self._poolSingleton[x]

    def getHash(self) -> int:
        return self._hash

    def __hash__(self) -> int:
        return self._hash


if __name__ == "__main__":
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    H1, H2 = FastHashSet(), FastHashSet()
    preHash1, preHash2 = [0], [0]
    for i in range(n):
        H1.add(nums1[i])
        preHash1.append(H1.getHash())
        H2.add(nums2[i])
        preHash2.append(H2.getHash())

    q = int(input())
    for _ in range(q):
        x, y = map(int, input().split())
        print("Yes" if preHash1[x] == preHash2[y] else "No")
