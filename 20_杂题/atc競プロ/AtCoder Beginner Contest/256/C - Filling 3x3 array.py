"""
3<=hi,wi<=30
分别表示行的和/列的和
每个格子填正整数一共有多少种可能的情况
只需要30^4枚举
"""
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def main() -> None:
    h1, h2, h3, w1, w2, w3 = map(int, input().split())
    res = 0
    for n1 in range(1, h1 + 1):
        for n2 in range(1, h1 + 1):
            n3 = h1 - n1 - n2
            if n3 <= 0:
                break
            for n4 in range(1, w1 + 1):
                for n5 in range(1, w2 + 1):
                    n6 = h2 - n4 - n5
                    if n6 <= 0:
                        break
                    n7 = w1 - n1 - n4
                    n8 = w2 - n2 - n5
                    if n7 <= 0 or n8 <= 0:
                        continue
                    cand1, cand2 = w3 - n3 - n6, h3 - n7 - n8
                    if cand1 != cand2 or cand1 <= 0:
                        continue
                    res += 1

    print(res)


if sys.argv[0].startswith(r"e:\test"):
    while True:
        try:
            main()
        except (EOFError, ValueError):
            break
else:
    main()
