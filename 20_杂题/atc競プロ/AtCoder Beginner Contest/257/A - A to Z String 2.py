import string
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def main() -> None:
    n, k = map(int, input().split())
    res = "".join([char * n for char in string.ascii_uppercase])
    print(res[k - 1])


while True:
    try:
        main()
    except (EOFError, ValueError):
        break
