# AtCoder Beginner Contest 259 - SGColin的文章 - 知乎
# https://zhuanlan.zhihu.com/p/539701972

# 给定n个数字的标准分解，将其中的某一个变成1，问操作后所有数字的最小公倍数有多少种不同的可能性?
# !结论是所有数字的最小公倍数等于每个质因数p的指数取所有数字对应质因数指数的max 。
# 一个数字变成1相当于对于LCM 什么都不提供，那么什么时候会导致LCM变化呢?
# !首先他的某一个质因数指数要和lcm对应的相同，其次这个最大值在所有数字中是唯一的。(拥有唯一的最大值)
# 因此开两个数组( unordered_map实现): mx[i]记录质因数i出现过的最大指数是多少，cnt[i]记录有多少个数字对应这个最大指数。
# 那么一个数字有贡献也就对应于e[i] = mx[i] && cnt[i] == 1 。
# !此外没有影响的所有数字总体会对答案产生一个贡献，即原本所有数的 LCM 。

from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n = int(input())
    nums = [defaultdict(int) for _ in range(n)]
    maxE = defaultdict(int)  # 每个质因数出现的最大指数
    maxCount = defaultdict(int)  # 每个质因数出现的最大指数的次数
    for i in range(n):
        m = int(input())
        for _ in range(m):
            p, e = map(int, input().split())
            nums[i][p] += e
            if e == maxE[p]:
                maxCount[p] += 1
            elif e > maxE[p]:
                maxE[p] = e
                maxCount[p] = 1

    res1, res2 = 0, 0  # 替换后变的情形，替换后不变的情形
    for num in nums:
        for p, e in num.items():
            if maxE[p] == e and maxCount[p] == 1:
                res1 += 1
                break
        else:
            res2 = 1

    print(res1 + res2)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
