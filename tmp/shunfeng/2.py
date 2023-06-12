from collections import deque

import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# 	顺丰小哥有个货物排成一排，第个货物的质量是，他准备取其中的一些货物。由于他不希望取完货后空隙过大，因此他希望取的货物中最多只有两件是相邻的，但是其余的货物都不能相邻。


# 	顺丰小哥突发奇想，他希望最终所有货物的质量乘积的结尾有尽可能多的0。你能帮他求出乘积末尾的0的最大数量吗？


#  输入描述 第一行输入一个正整数，代表货物的数量。第二行输入个正整数，用来表示每个货物的质量。 输出描述 一个整数，代表乘积末尾0的最大数量。
if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    C2, C5 = [0] * n, [0] * n
    for i, v in enumerate(nums):
        while v % 2 == 0:
            C2[i] += 1
            v //= 2
        while v % 5 == 0:
            C5[i] += 1
            v //= 5

    # @lru_cache(None)
    # def dfs(index: int, pre: bool, used: bool) -> Tuple[int, int]:
    #     if index >= n:
    #         return 0, 0
    #     # 不选
    #     res2, res5 = dfs(index + 1, False, used)
    #     # 选
    #     if not pre:
    #         tmp2, tmp5 = dfs(index + 1, True, False)
    #         cand1 = tmp2 + count2[index]
    #         cand2 = tmp5 + count5[index]
    #         if min(cand1, cand2) > min(res2, res5):
    #             res2, res5 = cand1, cand2

    #     elif not used:
    #         tmp2, tmp5 = dfs(index + 1, True, True)
    #         cand1 = tmp2 + count2[index]
    #         cand2 = tmp5 + count5[index]
    #         if min(cand1, cand2) > min(res2, res5):
    #             res2, res5 = cand1, cand2

    #     return res2, res5

    # res = dfs(0, False, False)
    # dfs.cache_clear()
    # print(min(res))

    #

    # !都不相邻
    c2, c5 = 0, 0
    # 所有奇数
    for i in range(0, n, 2):
        c2 += C2[i]
        c5 += C5[i]
    res1 = min(c2, c5)

    c2, c5 = 0, 0
    # 所有偶数
    for i in range(1, n, 2):
        c2 += C2[i]
        c5 += C5[i]
    res2 = min(c2, c5)

    # preMax = [0] * (n + 1)
    # max1, max2 = 0, 0
    # c21, c51 = 0, 0
    # c22, c52 = 0, 0
    # for i in range(n):
    #     if i % 2 == 0:
    #         c21 += C2[i]
    #         c51 += C5[i]
    #         max1 = max(c21, c51)
    #     else:
    #         c22 += C2[i]
    #         c52 += C5[i]
    #         max2 = max(c22, c52)
    #     preMax[i + 1] = max(max1, max2)

    # sufMax = [0] * (n + 1)
    # max1, max2 = 0, 0
    # c21, c51 = 0, 0
    # c22, c52 = 0, 0
    # for i in range(n - 1, -1, -1):
    #     if i % 2 == 0:
    #         c21 += C2[i]
    #         c51 += C5[i]
    #         max1 = max(c21, c51)
    #     else:
    #         c22 += C2[i]
    #         c52 += C5[i]
    #         max2 = max(c22, c52)
    #     sufMax[i] = max(max1, max2)

    # 前缀,所有的偶数位置
    preCount1 = [(0, 0) for _ in range(n + 1)]
    preCount2 = [(0, 0) for _ in range(n + 1)]
    c2, c5 = 0, 0
    d2, d5 = 0, 0
    for i in range(n):
        if i % 2 == 0:
            c2 += C2[i]
            c5 += C5[i]
        else:
            d2 += C2[i]
            d5 += C5[i]
        preCount1[i + 1] = (c2, c5)
        preCount2[i + 1] = (d2, d5)

    sufCount1 = [(0, 0) for _ in range(n + 1)]
    sufCount2 = [(0, 0) for _ in range(n + 1)]
    c2, c5 = 0, 0
    d2, d5 = 0, 0
    for i in range(n - 1, -1, -1):
        if i % 2 == 0:
            c2 += C2[i]
            c5 += C5[i]
        else:
            d2 += C2[i]
            d5 += C5[i]
        sufCount1[i] = (c2, c5)
        sufCount2[i] = (d2, d5)
    # print(preCount1, preCount2)
    # [(0, 0), (0, 2), (0, 2), (0, 5), (0, 5)] [(0, 0), (0, 0), (1, 1), (1, 1), (7, 1)]

    # print(C2, C5)
    # print(sufCount1, sufCount2)
    # [0, 1, 0, 6] [2, 1, 3, 0]
    # [(0, 5), (0, 3), (0, 3), (0, 0), (0, 0)] [(7, 1), (7, 1), (6, 0), (6, 0), (0, 0)]

    # 枚举相邻的元素位置???
    res = max(res1, res2)
    for i in range(n - 1):
        c2, c5 = C2[i] + C2[i + 1], C5[i] + C5[i + 1]
        # 左侧两种选择
        left = i - 1
        leftSelect = [(0, 0)]

        if left >= 0:
            leftC2, leftC5 = preCount1[left]
            leftSelect.append((leftC2, leftC5))
            leftC2, leftC5 = preCount2[left]
            leftSelect.append((leftC2, leftC5))
        right = i + 3
        rightSelect = [(0, 0)]
        if right <= n:
            rightC2, rightC5 = sufCount1[right]
            rightSelect.append((rightC2, rightC5))
            rightC2, rightC5 = sufCount2[right]
            rightSelect.append((rightC2, rightC5))

        cand = 0
        for leftC2, leftC5 in leftSelect:
            for rightC2, rightC5 in rightSelect:
                cand = max(cand, min(leftC2 + c2 + rightC2, leftC5 + c5 + rightC5))
        res = max(res, cand)

    print(res)
