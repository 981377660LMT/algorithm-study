import random
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 長さ N の文字列 S および整数 i (0≤i≤N) に対して、f
# i
# ​
#  (S) を、

# S の先頭 i 文字
# S を反転した文字列
# S の末尾 N−i 文字
# をこの順に連結した文字列と定義します。 例えば、S= abc、i=2 のとき、f
# i
# ​
#  (S)= abcbac です。

# 長さ 2N の文字列 T が与えられます。 f
# i
# ​
#  (S)=T を満たす長さ N の文字列 S と整数 i (0≤i≤N) の組を 1 つ見つけてください。 そのような S,i の組が存在しない場合は、それを報告してください。
if __name__ == "__main__":

    # 从移除s[::-1]后的字符串中取
    def getChar1(insertPos: int, i: int) -> str:
        if i < insertPos or i >= n + insertPos:
            return t[i]
        return t[n + i]

    # 从s[::-1]中取
    def getChar2(insertPos: int, i: int) -> str:
        end = n + insertPos
        return t[end - i - 1]

    n = int(input())
    t = input()

    fail = 0
    for i in range(n + 1):  # insert position
        # 随机检查50个下标
        for _ in range(50):
            j = random.randint(0, n - 1)
            if getChar1(i, j) != getChar2(i, j):
                break
        else:
            s1 = t[:i] + t[n + i :]
            s2 = t[i : n + i]
            if s1 == s2[::-1]:
                print(s1)
                print(i)
                exit(0)
            else:
                fail += 1
        if fail >= 50:
            print(-1)
            exit(0)
    print(-1)
