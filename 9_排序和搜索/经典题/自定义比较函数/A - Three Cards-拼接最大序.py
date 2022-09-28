# 拼接最大数
# n个数中选取三个数 任意顺序拼接
# !求能拼接的最大数
# n<=2e5 ai<=1e6

# !n个数选择三个数使得拼接出来的整数最大
# !1.选择的整数越大越好
# !2.拼接最大序
# 反例: 98 987 23 12

from functools import cmp_to_key
from heapq import nlargest


def mergeSort(s1: str, s2: str) -> int:
    """拼接两个字符串，字典序最小"""
    return -1 if s1 + s2 < s2 + s1 else 1


n = int(input())
nums = list(map(int, input().split()))

max3 = nlargest(3, nums)
max3.sort(key=cmp_to_key(lambda s1, s2: mergeSort(str(s1), str(s2))), reverse=True)
print(*max3, sep="")
