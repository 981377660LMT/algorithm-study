# FastHash, 集合哈希

from collections import defaultdict
from random import randint
from typing import Any, Iterable, List, Optional


class FastHashSet:
    """可以快速计算哈希值的集合."""

    _poolSingleton = defaultdict(lambda: randint(1, (1 << 61) - 1))

    __slots__ = ("_set", "_hash")

    def __init__(self, iterable: Optional[Iterable[Any]] = None) -> None:
        self._set = set()
        self._hash = 0
        if iterable is not None:
            for x in iterable:
                self.add(x)

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

    def symmetricDifference(self, other: "FastHashSet") -> int:
        return self._hash ^ other._hash

    def clear(self) -> None:
        self._hash = 0
        self._set.clear()

    def __hash__(self) -> int:
        return self._hash


class FastHashCounter:
    """可以快速计算哈希值的Counter."""

    _poolSingleton = defaultdict(lambda: randint(1, (1 << 61) - 1))

    __slots__ = "_hash", "_counter"

    def __init__(self) -> None:
        self._counter = defaultdict(int)
        self._hash = 0

    def add(self, x: int) -> None:
        self._counter[x] += 1
        self._hash += self._poolSingleton[x]

    def discard(self, x: int) -> bool:
        if self._counter[x] == 0:
            return False
        self._counter[x] -= 1
        self._hash -= self._poolSingleton[x]
        return True

    def getHash(self) -> int:
        return self._hash

    def symmetricDifference(self, other: "FastHashCounter") -> int:
        return self._hash ^ other._hash

    def clear(self) -> None:
        self._hash = 0
        self._counter.clear()

    def __hash__(self) -> int:
        return self._hash


class FastHashChessBoard:
    """
    可以快速计算哈希值的棋盘.
    棋子的(id,位置)唯一确定一个`棋盘上的`棋子.
    所有棋子的哈希值唯一确定棋盘的哈希值.
    """

    ChessId = int  # 棋子的(id,位置)唯一确定一个棋盘上的棋子 => 6x2x64 = 768

    _poolSingleton = defaultdict(lambda: randint(1, (1 << 61) - 1))

    __slots__ = ("_hash", "_board")

    def __init__(self, chess: List[ChessId], initPos: List[int]) -> None:
        """指定每个棋子的编号和初始位置, 初始化棋盘."""
        self._hash = 0
        self._board = defaultdict(int)
        for chessId, pos in zip(chess, initPos):
            self._board[chessId] = pos
            self._hash ^= self._poolSingleton[(chessId, pos)]

    def move(self, chessId: ChessId, toPos: int) -> None:
        """将编号为chessId的棋子移动到toPos位置."""
        oldPos = self._board[chessId]
        newPos = toPos
        oldHash = self._poolSingleton[(chessId, oldPos)]
        newHash = self._poolSingleton[(chessId, newPos)]
        self._hash ^= oldHash ^ newHash
        self._board[chessId] = newPos

    def getHash(self) -> int:
        return self._hash

    def symmetricDifference(self, other: "FastHashChessBoard") -> int:
        return self._hash ^ other._hash

    def __hash__(self) -> int:
        return self._hash

    def __repr__(self) -> str:
        return repr(self._board)


class FastHashRange:
    """可以快速计算哈希值的区间."""

    _poolSingleton = defaultdict(lambda: randint(1, (1 << 61) - 1))

    __slots__ = "_hash"

    def __init__(self) -> None:
        self._hash = 0

    def add(self, left: int, right: int, delta: int) -> None:
        """区间[left,right]中每个数加上delta.
        0 <= left <= right < n.
        """
        self._hash += (self._poolSingleton[right] - self._poolSingleton[left]) * delta

    def getHash(self) -> int:
        return self._hash

    def symmetricDifference(self, other: "FastHashRange") -> int:
        return self._hash ^ other._hash

    def __hash__(self) -> int:
        return self._hash


class AllCountMultipleOfKChecker:
    """判断数据结构中每个数出现的次数是否均k的倍数."""

    _poolSingleton = defaultdict(lambda: randint(1, (1 << 61) - 1))

    __slots__ = ("_hash", "_modCounter", "_k")

    def __init__(self, k: int) -> None:
        self._hash = 0
        self._modCounter = defaultdict(int)
        self._k = k

    def add(self, x: int) -> None:
        count, random = self._modCounter[x], self._poolSingleton[x]
        self._hash -= count * random
        count += 1
        if count == self._k:
            count = 0
        self._hash += count * random
        self._modCounter[x] = count

    def remove(self, x: int) -> None:
        """删除前需要保证x在集合中存在."""
        count, random = self._modCounter[x], self._poolSingleton[x]
        self._hash -= count * random
        count -= 1
        if count == -1:
            count = self._k - 1
        self._hash += count * random
        self._modCounter[x] = count

    def query(self) -> bool:
        return self._hash == 0

    def getHash(self) -> int:
        return self._hash

    def clear(self) -> None:
        self._hash = 0
        self._modCounter.clear()
