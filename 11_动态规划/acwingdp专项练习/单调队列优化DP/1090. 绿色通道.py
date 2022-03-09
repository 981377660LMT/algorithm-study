# 二分+单调队列dp
# 高二数学《绿色通道》总共有 n 道题目要抄，编号 1,2,…,n，抄第 i 题要花 ai 分钟。
# 小 Y 决定只用不超过 t 分钟抄这个，因此必然有空着的题。
# 每道题要么不写，要么抄完，不能写一半。
# 下标连续的一些空题称为一个空题段，它的长度就是所包含的题目数。
# 这样应付自然会引起马老师的愤怒，最长的空题段越长，马老师越生气。
# 现在，小 Y 想知道他在这 t 分钟内写哪些题，才能够尽量减轻马老师的怒火。
# 由于小 Y 很聪明，你只要告诉他最长的空题段至少有多长就可以了，不需输出方案。

'''
二分搜索可能的空题段的长度，用单调队列优化的DP验证此
空题段长度下，能够找到合法的选择方案
'''
from collections import deque

n, m = map(int, input().split())
nums = list(map(int, input().split()))


def check(mid: int) -> bool:
    ...


left, right = 0, n
while left <= right:
    mid = (left + right) >> 1
    if check(mid):
        right = mid - 1
    else:
        left = mid + 1
print(left)
