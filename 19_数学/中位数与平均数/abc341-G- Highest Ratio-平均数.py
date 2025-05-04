# G - Highest Ratio
# https://atcoder.jp/contests/abc341/tasks/abc341_g
# !给定一个长度为N的数组，对于每个左端点开始的数组，求最大子段平均值.
#
# !类似 1792. 最大平均通过率
# https://leetcode.cn/problems/maximum-average-pass-ratio/solution/zui-da-ping-jun-tong-guo-lu-by-leetcode-dm7y3/
#
# !构建一个下凸壳（lower convex hull），只保留那些可能成为最优解的点。

from itertools import accumulate


if __name__ == "__main__":
    N = int(input())
    A = list(map(int, input().split()))
    presum = [0] + list(accumulate(A))

    res = [0.0] * N
    stack = [(0, 1)]

    for i in range(N - 1, -1, -1):
        top = A[i]
        bottom = 1
        while stack[-1][0] / stack[-1][1] > top / bottom:
            top += stack[-1][0]
            bottom += stack[-1][1]
            stack.pop()
        stack.append((top, bottom))
        res[i] = top / bottom

    print(*res, sep="\n")
