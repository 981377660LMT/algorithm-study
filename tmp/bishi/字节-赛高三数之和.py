# 赛高的三数之和

# D有一个长度为n的数组,现在他也遇到了经典的三数和问题,
# 但是这个问题有些不一样，请你帮助他。
# !问题描述是这样的,如果在数组a，对任意的i、j、k都能找到l，
# !使得ai +aj+ak =al,那么这个数组就是被称为三数和赛高数组,
# 特别的i,j, k,I需要满足(1<=ijsk<=n,1<=<=n)。
# n<=2e5

# from random import sample
# import string


# T = int(input())
# for _ in range(T):
#     n = int(input())
#     nums = list(map(int, input().split()))
#     s = set(nums)
#     for _ in range(int(1e5)):
#         a, b, c = sample(nums, 3)
#         if a + b + c not in s:
#             print("No")
#             break
#     print("Yes")

#######################################################################################
# 有一个小写字母组成的字符串str，字符串长度为n;
# !字符串中`某些位置`可以执行修复操作:将这个位置的字符，替换为a~z中的一个字符;
# !这个修复操作最多可以执行m次;
# !现在想知道修复之后,字符串str由相同字符组成的子串最大长度是多少。
# n,m<=2000
# 26*2000*log(2000)

import string


INF = int(1e9)
T = int(input())
for _ in range(T):
    n, m = map(int, input().split())
    s = input()
    canModify = input()

    # 枚举哪一种字符
    res = 1
    for char in string.ascii_lowercase:

        def check(mid: int) -> bool:
            """修改m次,能否构造char组成的长度为mid的相同字符子串(定长滑窗)"""
            res, cost = INF, 0
            for right in range(n):
                if s[right] != char:
                    cost += INF if canModify[right] == "0" else 1
                if right >= mid:
                    if s[right - mid] != char:
                        cost -= INF if canModify[right - mid] == "0" else 1
                if right >= mid - 1:
                    res = min(res, cost)
            return res <= m

        left, right = 1, n
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        res = max(res, right)

    print(res)
