# FastHash

from collections import defaultdict
from random import randint
from typing import List


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

    def symmetricDifference(self, other: "FastHashSet") -> int:
        return self._hash ^ other._hash

    def __hash__(self) -> int:
        return self._hash


class FastHashCounter:
    """快速计算哈希值的Counter."""

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

    def __hash__(self) -> int:
        return self._hash


class FastHashChessBoard:
    """
    快速计算哈希值的棋盘.
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
