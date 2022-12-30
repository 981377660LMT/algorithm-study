from random import randint
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 行 M 列のグリッドがあります． あなたはグリッドの各マスに 0 以上 2
# K
#  −1 以下の整数を書き込み，以下の条件を満たしたいです．

# 左上のマスを出発し，右または下に隣接するマスへの移動を繰り返して，右下のマスへと至るパスを考える． ここで，通ったマス (始点終点を含む) に書かれた整数の総 XOR が 0 になるパスを，よいパスと呼ぶことにする．
# よいパスはちょうど 1 つだけ存在し，それは文字列 S が表すパスである． 文字列 S が表すパスとは，各 i (1≤i≤N+M−2) について，i 回目の移動の際，S の i 文字目が R なら右，D なら下に進むようなパスである．
# 条件を満たす整数の書き込み方が存在するかどうか判定してください．
# 1 つの入力ファイルにつき，T 個のテストケースを解いてください．

# 1≤T≤100
# 2≤N,M≤30
# 1≤K≤30
# S はちょうど N−1 個の D と M−1 個の R からなる文字列
# 入力される数はすべて整数

# 各ケースは以下の形式で与えられる．

# N M K
# S

# 各ケースに対し，条件を満たす整数の書き込み方が存在する場合は Yes を，存在しないならば No を出力せよ．
# 出力中の各文字は英大文字・小文字のいずれでもよい．

# 折半枚举
if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        ROW, COL, K = map(int, input().split())
        S = input()
        ok = False
        # 随机填充
        for _ in range(30):
            grid = [[randint(0, 2**K - 1) for _ in range(COL)] for _ in range(ROW)]
            curXor = grid[0][0]
            r, c = 0, 0
            for i in range(ROW + COL - 2):
                if S[i] == "R":
                    c += 1
                else:
                    r += 1
                curXor ^= grid[r][c]
            if curXor == 0:
                ok = True
                break

        print("Yes" if ok else "No")
