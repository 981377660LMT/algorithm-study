# 题目背景
# 在《LIAR GAME》中，小 E 看到了一个有趣的游戏。

# 题目描述
# 这个游戏名叫《走私游戏》。游戏规则大概是这样的：一名玩家扮演走私者，一名玩家扮演检察官。走私者可以将
# �
# x 日元（
# �
# x 为
# [
# 0
# ,
# �
# ]
# [0,n] 内的整数，由走私者决定）秘密放入箱子中，而检查官需要猜测箱子中的金额，猜测范围也为
# [
# 0
# ,
# �
# ]
# [0,n] 中的整数。假设检察官猜了
# �
# y。如果
# �
# =
# �
# x=y，则走私失败，走私者一分钱也拿不到。如果
# �
# >
# �
# x>y，则走私成功，走私者可以从检查官那里拿走
# �
# x 日元。如果
# �
# <
# �
# x<y，则走私失败，但是由于冤枉检察官需要赔付给走私者
# �
# /
# 2
# y/2 日元。游戏分有限回合进行。双方轮流做走私者和检察官。

# 可以证明，最优情况下每个回合走私者会采用同一种策略，检察官也会采用同一种策略。小 E 想知道在一个回合中，双方的最优策略分别是什么。

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    MOD = 998244353
    n = int(input()) + 1
    fac = [1] * (n + 1)
    ifac = [1] * (n + 1)
    for i in range(1, n + 1):
        fac[i] = fac[i - 1] * i % MOD
    ifac[n] = pow(fac[n], MOD - 2, MOD)
    for i in range(n, 0, -1):
        ifac[i - 1] = ifac[i] * i % MOD

    q = [1]
    acc = [0, 1]
    for i in range(1, n):
        q.append(((i * q[-1] * ifac[2] + acc[-2] * ifac[i]) * ifac[i]) % MOD)
        acc.append((acc[-1] + q[-1]) % MOD)

    tmp = acc[-1]
    v = pow(tmp, MOD - 2, MOD)
    q = [num * v % MOD for num in q]
    print(*q)

    q = [1]
    acc = [0, 1]
    for i in range(1, n):
        q.append(((i * q[-1] + acc[-2]) * ifac[i] * 2) % MOD)
        acc.append((acc[-1] + q[-1]) % MOD)

    tmp = acc[-1]
    v = pow(tmp, MOD - 2, MOD)
    q = [num * v % MOD for num in q]
    print(*q)
