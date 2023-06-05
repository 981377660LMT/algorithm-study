# https://nyaannyaan.github.io/library/game/partisan-game.hpp
# PartisanGame-非公平博弈/非不偏ゲーム

# TODO

from typing import Callable, Generic, List, Tuple, TypeVar
from SurrealNumber import SurrealNumber


Game = TypeVar("Game")


class PartisanGame(Generic["Game"]):
    __slots__ = ("_mp", "_f")

    def __init__(self, f: Callable[[Game], Tuple[List[Game], List[Game]]]):
        self._mp = dict()
        self._f = f

    def zero(self) -> "SurrealNumber":
        return SurrealNumber(0)

    def get(self, g: "Game") -> "SurrealNumber":
        res = self._mp.get(g, None)
        if res is not None:
            return res
        res = self._get(g)
        self._mp[g] = res
        return res

    def _get(self, g: "Game") -> "SurrealNumber":
        gl, gr = self._f(g)
        if not gl and not gr:
            return self.zero()
        ls = [self.get(cg) for cg in gl]
        rs = [self.get(cg) for cg in gr]
        sl = max(ls) if ls else self.zero()
        sr = min(rs) if rs else self.zero()
        if not rs:
            return sl.larger()
        if not ls:
            return sr.smaller()
        # assert sl < sr
        return SurrealNumber.reduce(sl, sr)
