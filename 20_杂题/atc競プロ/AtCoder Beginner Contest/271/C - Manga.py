"""
现在有一系列书, 编号从1 到 1e9, 高桥持有 n 本书, 第 i 本书的编号为 ai , 

在高桥开始阅读之前, 可以进行以下操作:
如果手中少于2本书, 则什么也不干, 否则就选择两本书卖出去, 换取指定的一本书

操作完成之后, 
高桥将按顺序阅读, 先阅读编号为1的书, 然后是编号为2的书.... 直到连接不上下一本书为止, 
之后的书高桥都不会阅读, 现在问高桥能阅读到的最大的书的编号为多少?

!二分答案比较方便
"""

from collections import Counter
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    counter = Counter(nums)

    def check(mid: int) -> bool:
        money = 0
        for key in counter:
            if key > mid:
                money += counter[key]
            else:
                money += counter[key] - 1

        for i in range(1, mid + 1):
            if counter[i] == 0:
                if money >= 2:
                    money -= 2
                else:
                    return False

        return True

    left, right = 1, int(1e10)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1

    print(right)
