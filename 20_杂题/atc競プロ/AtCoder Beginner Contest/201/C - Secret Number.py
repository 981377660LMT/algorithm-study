import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 高橋くんは、暗証番号を忘れてしまいました。暗証番号は 0 から 9 までの数字のみからなる 4 桁の文字列で、0 から始まる場合もあります。

# 0 から 9 までの各数字について、高橋くんは以下のように記憶しています。彼の記憶は長さ 10 の文字列 S
# 0
# ​
#  S
# 1
# ​
#  …S
# 9
# ​
#   によって表されます。

# S
# i
# ​
#   が o のとき : 数字 i は暗証番号に確実に含まれていた。
# S
# i
# ​
#   が x のとき : 数字 i は暗証番号に確実に含まれていなかった。
# S
# i
# ​
#   が ? のとき : 数字 i が暗証番号に含まれているか分からない。
# 高橋くんが忘れてしまった暗証番号としてあり得るものは何通りありますか？
if __name__ == "__main__":

    def dfs(index: int, path: str) -> None:
        global res
        if index == 4:
            res += all(c in path for c in must)
            return
        for cur in must + may:
            dfs(index + 1, path + cur)

    s = input()
    must = [str(i) for i in range(10) if s[i] == "o"]
    may = [str(i) for i in range(10) if s[i] == "?"]
    res = 0
    dfs(0, "")
    print(res)
