# 主人公N岁的时候身高是T，已知他[1,X]期间每年长D，后面不长个子，问M岁的时候身高多少

# 先求0岁时的身高 再求m岁的身高

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    curAge, queriedAge, growLimit, curHeight, delta = map(int, input().split())
    zero = curHeight - min(curAge, growLimit) * delta
    print(zero + min(queriedAge, growLimit) * delta)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
