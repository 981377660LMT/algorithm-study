# k张优惠券 每个降价x (可叠加,但不能降至0以下)
# 求买下所有商品的最小花费

# 要让优惠券最多的方案
# !如果优惠券能减去商品价格大于x的话,就贪心地用优惠券(-x)
# !如果所有的商品剩下的价格都小于x元了之后再还有优惠券剩下,
# !剩下的每个商品费用都是小于x的，我们从大到小贪心即可。

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, k, x = map(int, input().split())
    prices = list(map(int, input().split()))

    remain = k
    for i in range(n):
        count = prices[i] // x
        cur = min(remain, count)
        remain -= cur
        prices[i] -= cur * x

    if remain > 0:
        prices.sort(reverse=True)
        for i in range(n):
            if remain == 0:
                break
            prices[i] = 0
            remain -= 1

    print(sum(prices))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
