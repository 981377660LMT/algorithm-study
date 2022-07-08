# 高橋君は「A 秒間秒速 B メートルで歩き、C 秒間休む」ことを繰り返します。
# 青木君は「D 秒間秒速 E メートルで歩き、F 秒間休む」ことを繰り返します。
# 二人が同時にジョギングを始めてから X 秒後、高橋君と青木君のうちどちらが長い距離を進んでいますか？

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    def cal(goTime: int, speed: int, stopTime: int, allTime: int) -> int:
        res = 0
        while allTime > 0:
            go = min(allTime, goTime)
            res += go * speed
            allTime -= go
            allTime -= stopTime
        return res

    a, b, c, d, e, f, x = map(int, input().split())
    res1, res2 = cal(a, b, c, x), cal(d, e, f, x)
    if res1 > res2:
        print("Takahashi")
    elif res1 < res2:
        print("Aoki")
    else:
        print("Draw")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
