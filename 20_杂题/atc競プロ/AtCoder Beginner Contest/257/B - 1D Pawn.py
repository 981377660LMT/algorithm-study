# 左から L i番目のコマが一番右のマスにあるならば何も行わない。
# そうでない時、左から Li番目のコマがあるマスの 1 つ右のマスにコマが無いならば、
# 左から L i番目のコマを 1 つ右のマスに移動させる。
# 1 つ右のマスにコマがあるならば、何も行わない。

import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def main() -> None:
    n, _, _ = map(int, input().split())
    pos = list(map(int, input().split()))
    Q = [int(num) - 1 for num in input().split()]  # !虽然这样命名不太好 但Q比queries好写

    for i in Q:
        if (pos[i] == n) or (i + 1 < len(pos) and pos[i] + 1 == pos[i + 1]):
            continue
        pos[i] += 1
    print(*pos)


while True:
    try:
        main()
    except (EOFError, ValueError):
        break
