from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 高橋君は
# N 匹のモンスターに順に出会います。
# i 匹目
# (1≤i≤N) のモンスターの強さは
# A
# i
# ​
#   です。

# 高橋君はそれぞれのモンスターについて逃がすか倒すか選ぶことができます。
# 高橋君はそれぞれの行動によって次のように経験値を得ます。

# モンスターを逃がした場合、得られる経験値は
# 0 である。
# 強さが
# X のモンスターを倒したとき、経験値を
# X 得る。
# ただし、モンスターを倒すのが偶数回目（
# 2,
# 4,
# … 回目）であるとき、さらに追加で経験値を
# X 得る。
# N 匹から高橋君が得た経験値の合計としてあり得る最大の値を求めてください。
if __name__ == "__main__":
    N = int(input())
    A = list(map(int, input().split()))
