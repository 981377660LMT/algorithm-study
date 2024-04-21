# from titan_pylib.data_structures.safe_hash.hash_dict import HashDict
from typing import Dict, Any
import random


class HashDict:
    _xor = random.randrange(10000000, 1000000000)

    def __init__(self):
        self._data: Dict[int, Any] = {}

    def __setitem__(self, key: int, val: Any):
        self._data[key ^ HashDict._xor] = val

    def __getitem__(self, key: int) -> Any:
        return self._data[key ^ HashDict._xor]

    def __delitem__(self, key: int):
        del self._data[key ^ HashDict._xor]

    def __contains__(self, key: int):
        return key ^ HashDict._xor in self._data

    def __len__(self):
        return len(self._data)

    def keys(self):
        return (k ^ HashDict._xor for k in self._data.keys())

    def values(self):
        return (v for v in self._data.values())

    def items(self):
        return ((k ^ HashDict._xor, v) for k, v in self._data.items())

    def __str__(self):
        return "{" + ", ".join(f"{k}: {v}" for k, v in self.items()) + "}"
