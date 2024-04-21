# from titan_pylib.data_structures.safe_hash.hash_set import HashSet
from typing import Iterable, Set
import random


class HashSet:
    _xor = random.randrange(10000000, 1000000000)

    def __init__(self, a: Iterable[int] = []):
        self._data: Set[int] = set(x ^ HashSet._xor for x in a)

    def add(self, key: int) -> None:
        self._data.add(key ^ HashSet._xor)

    def discard(self, key: int) -> None:
        self._data.discard(key ^ HashSet._xor)

    def remove(self, key: int) -> None:
        self._data.remove(key ^ HashSet._xor)

    def __contains__(self, key: int):
        return key ^ HashSet._xor in self._data

    def __len__(self):
        return len(self._data)

    def __iter__(self):
        return (k ^ HashSet._xor for k in self._data.__iter__())

    def __str__(self):
        return "{" + ", ".join(sorted(map(str, self))) + "}"

    def __repr__(self):
        return f"{self.__class__.__name__}({self})"
