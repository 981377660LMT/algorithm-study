# # !C - トーナメント
# # k<=10 Ri<=4000
# # !2^k个人相邻两人对战 求第i个人最后胜出的概率
# !dp[i][j] 表示第i轮j胜出的概率
# 第0轮看1个人 第1轮看2个人 第三轮看4个人 。。。
# https://atcoder.jp/contests/tdpc/submissions/33741579

# !dp + 完全二叉树
# 最佳运动员的比拼回合2
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


k = int(input())
n = 1 << k
ratings = [int(input()) for _ in range(2**k)]


def cal(i: int, j: int) -> float:
    """i与j pk i胜出的概率"""
    p1 = 1 / (1 + pow(10, (ratings[j] - ratings[i]) / 400))
    return p1


def check(a: int, b: int, i: int) -> bool:
    """a和b是否可能在第i轮相遇

    画二叉树验证发现子节点与父结点关系 第i轮相遇必须满足两个条件
    1. a >> (i - 1)  与 b >> (i - 1)  不能相等(即作为对手位置)
    2. (a >> i) 与 (b >> i) 相等(即恰好相遇)
    """
    if a >> (i - 1) == b >> (i - 1):
        return False
    return a >> i == b >> i


dp = [1.0] * n
for i in range(1, k + 1):  # !每轮要淘汰一半人 完全二叉树
    ndp = [0.0] * n
    for pre in range(n):
        for cur in range(n):
            if check(cur, pre, i):  # !pre和cur是否可能在第i轮相遇
                ndp[cur] += dp[cur] * dp[pre] * cal(cur, pre)  # !之前他们都赢*这次cur赢
    dp = ndp


for num in dp:
    print(f"{num:.15f}")
