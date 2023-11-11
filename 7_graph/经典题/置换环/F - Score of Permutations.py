# 题意:对于序列P=(P1,P2，…,PN)，
# 按照如下规则定义序列P的分数S(P):
# 有N个人，分别编号为1,2,….,N。第i个人有一个编号为i的球。
# 每当Snuke喊一次，编号为i的人都会同时把自己的球给编号为pi的人。
# 当每个人的编号都和自己获得的球的编号相同时，Snuke停止叫喊。
# 分数为Snuke 喊的次数，保证分数是一个有限的值。
# !现在给出两个数字n,k(n≤50,k<=1e4)，计算所有的排列的分数的k次方的和

# !看到这个k次幂的就知道不会有什么神奇的性质了，得老老实实把每个答案都求出来了。

# !置换环
# 1.对于每一个置换环，每个球都恢复到原来的位置所需的次数为所有置换环长度的最小公倍数。
# 2.我们可以通过加入一个环来进行我们的dp
# !dp[使用了i个数][环的lcm为j] = 方案数
# 我们转移的时候可以可以枚举k(需要从剩下的n-i个元素中选取k个组成环)
# 为了不重复计数 每次枚举的k个中都要包含剩下的n-i个元素中的`最小值` 例如(1,2)->(3,4) 与(3,4)->(1,2)是相同的
# !dp[i][lcm(j, k)] = ∑dp[i][j] * C(n - i -1 , k - 1) * fac[k - 1]


from math import gcd
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


fac = [1]
ifac = [1]
for i in range(1, int(1e4) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


if __name__ == "__main__":
    n, K = map(int, input().split())
    dp = [dict() for _ in range(n + 1)]  # dp[i][lcm] = 方案数
    dp[0][1] = 1
    for i in range(n):
        for preLcm, preCount in dp[i].items():
            for k in range(1, n - i + 1):  # 枚举新加入环的长度k 每次选择环都要包含剩下的最小值
                curLcm = preLcm * k // gcd(preLcm, k)
                dp[i + k][curLcm] = (
                    dp[i + k].get(curLcm, 0) + preCount * C(n - i - 1, k - 1) * fac[k - 1]
                ) % MOD

    res = 0

    for lcm, count in dp[n].items():
        res = (res + count * pow(lcm, K, MOD)) % MOD
    print(res)
