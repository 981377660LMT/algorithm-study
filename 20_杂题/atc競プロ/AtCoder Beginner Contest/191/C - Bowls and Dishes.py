from itertools import product
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")

# 1,2,…,N の番号がついた N 個の皿と、
# 1,2,…,M の番号がついた M 個の条件があります。
# 条件 i は、皿 Aiと皿 Biの両方にボールが (1 個以上) 置かれているとき満たされます。
# 1,2,…,K の番号がついた K 人の人がいて、
# 人 i は皿 Ciか皿 Diのどちらか一方にボールを置きます。
# 満たされる条件の個数は最大でいくつでしょうか？
# k<=16

# !笛卡尔积
if __name__ == "__main__":
    n, m = map(int, input().split())
    pairs = [list(map(int, input().split())) for _ in range(m)]
    k = int(input())
    balls = [list(map(int, input().split())) for _ in range(k)]

    res = 0
    for select in product(*balls):
        S = set(select)
        cur = sum(((a in S) and (b in S)) for a, b in pairs)
        res = max(res, cur)
    print(res)
