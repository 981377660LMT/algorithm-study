from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    Sx, Sy = map(int, input().split())
    Tx, Ty = map(int, input().split())

    def getBlockXId(x: int, y: int) -> int:
        return x + ((y & 1) == 0)

    if Sx > Tx:
        Sx, Tx = Tx, Sx
    diffX = abs(getBlockXId(Tx, Ty) - getBlockXId(Sx, Sy))
    diffY = abs(Ty - Sy)
    print(max(diffX // 2, diffY))
