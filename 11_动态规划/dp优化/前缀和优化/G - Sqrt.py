"""
G给序列,序列里面只有一个数x(x<=9e18)。
假设序列最后一个数为t,你可以接一个属于1~t^(1/2)的整数在t后面,(加一个不超过根号x的数)
这样的操作无数次(1e100),问可以形成多少种不同的序列?
例如x为16时有5种
(16,4,2,1,1,1,…)
(16,4,1,1,1,1,…)
(16,3,1,1,1,1,…)
(16,2,1,1,1,1,…)
(16,1,1,1,1,1,…)

非常显然的 dp
前缀和优化dp + 根号分块计算
然后每个dpSum[x]在区间[t^2,t^2+2t]内的数字取根号是一样的
"""

# !公式推导变形优化
# !dp[v]=dpsum[v^(1/2)]=dp[1]+dp[2]+...+dp[v^(1/2)] 即
# !dp[v]=dpsum[1]+dpsum[2^(1/2)]+...+dpsum[v^(1/2)^(1/2)]
# !在每一个区间块 [t^2, t^2+2t] 内计算 `dp[t]*区间长度`


from collections import defaultdict
from math import isqrt
import sys
import os


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

# sqrt方法精度就会下降导致开根计算出错
# isqrt = lambda x: int(x**0.5)

# !isqrt可以避免精度问题 (采用牛顿迭代法计算开根)
# 一般的开根使用浮点数
# !当浮点数超过2^53-1, 转整数之后会出现精度问题


def main() -> None:
    # !1. 预处理dpSum
    dp = defaultdict(int, {1: 1})
    dpSum = defaultdict(int, {1: 1})
    for i in range(2, int(1e5)):
        dp[i] = dpSum[isqrt(i)]
        dpSum[i] = dpSum[i - 1] + dp[i]

    T = int(input())  # T个查询 T<=20
    for _ in range(T):
        n = int(input())
        sqrt2 = isqrt(n)

        res, remain, bid = 0, sqrt2, 1  # !bid表示区间块索引
        while remain:
            lower, upper = bid * bid, bid * bid + 2 * bid
            count = min(remain, upper - lower + 1)
            res += dpSum[bid] * count
            remain -= count
            bid += 1

        print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
