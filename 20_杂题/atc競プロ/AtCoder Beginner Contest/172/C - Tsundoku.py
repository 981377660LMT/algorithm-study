# C - Tsundoku
# !積読(买书不读书)
# 两个书桌：A和B。A书桌上有N本书，B书桌上有M本书。
# 读A书桌上的第i本书需要a[i]分钟，读B书桌上的第i本书需要b[i]分钟。
# 选择一个书桌，从最上面的书开始读，读完就移除
# （至于是A的开头还是B的开头，无所谓，可以交叉）。

# 最多给你k分钟，按照这样的规则拿书，最多可以拿多少本书？

# n,m<=2e5

# !两个栈,前缀和+枚举+二分
from bisect import bisect_right
from itertools import accumulate
from typing import List


def tsundoku(cost1: List[int], cost2: List[int]) -> int:
    preSum1 = [0] + list(accumulate(cost1))
    preSum2 = [0] + list(accumulate(cost2))
    res = 0
    for count1 in range(len(cost1) + 1):
        sum1 = preSum1[count1]
        if sum1 > k:
            break
        count2 = bisect_right(preSum2, k - sum1) - 1
        res = max(res, count1 + count2)
    return res


if __name__ == "__main__":
    n, m, k = map(int, input().split())
    cost1 = list(map(int, input().split()))
    cost2 = list(map(int, input().split()))
    print(tsundoku(cost1, cost2))
