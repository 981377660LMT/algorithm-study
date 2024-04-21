# from titan_pylib.data_structures.safe_hash.hash_defaultdict import HashDefaultdict
from collections import defaultdict
import random


class HashDefaultdict:
    _xor = random.randrange(10000000, 1000000000)

    def __init__(self, missing):
        self._data = defaultdict(missing)

    def __iter__(self):
        for k, v in self.items():
            yield k, v

    def __setitem__(self, key: int, value):
        self._data[key ^ HashDefaultdict._xor] = value

    def __getitem__(self, key: int):
        return self._data[key ^ HashDefaultdict._xor]

    def __delitem__(self, key: int):
        del self._data[key ^ HashDefaultdict._xor]

    def __contains__(self, item):
        return item ^ HashDefaultdict._xor in self._data

    def __len__(self):
        return len(self._data)

    def keys(self):
        return (k ^ HashDefaultdict._xor for k in self._data.keys())

    def values(self):
        return (v for v in self._data.values())

    def items(self):
        return ((k ^ HashDefaultdict._xor, v) for k, v in self._data.items())

    def __str__(self):
        return "HashDefaultdict({" + ", ".join(f"{k}: {v}" for k, v in self.items()) + "})"
