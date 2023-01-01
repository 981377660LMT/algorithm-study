class safedict(dict):
    def __init__(self):
        import random

        self._r = random.randrange(1 << 63)

    def __contains__(self, key):
        return super().__contains__(key ^ self._r)

    def __getitem__(self, key):
        return super().__getitem__(key ^ self._r)

    def __setitem__(self, key, value):
        return super().__setitem__(key ^ self._r, value)

    def get(self, key, default=None):
        return super().get(key ^ self._r, default)


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    q = int(input())
    mp = safedict()
    for _ in range(q):
        op, *args = map(int, input().split())
        if op == 0:
            k, v = args
            mp[k] = v
        elif op == 1:
            k = args[0]
            print(mp.get(k, 0))
