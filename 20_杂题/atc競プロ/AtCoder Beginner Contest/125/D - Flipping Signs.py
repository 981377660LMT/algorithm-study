# 每次操作可以选择数组相邻元素乘上-1
# 进行任意次操作,使得数组和最大

# 观察性质 统计负数个数的奇偶性
# !最后至多1个负数

# 1.如果负数偶数个,那么最后没有负数
# 2.如果负数奇数个,但是有0,那么最后没有负数
# !3.如果负数奇数个,没有0,那么最后有一个负数 让所有正负数里绝对值最小的数变为负数


import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    n = int(input())
    nums = list(map(int, input().split()))
    pos, zero, neg = [], [], []  # !保存绝对值
    for num in nums:
        if num > 0:
            pos.append(num)
        elif num == 0:
            zero.append(num)
        else:
            neg.append(-num)

    if len(neg) % 2 == 0 or len(zero) > 0:
        print(sum(pos) + sum(zero) + sum(neg))
    else:
        min_ = min(min(pos, default=INF), min(neg, default=INF))
        print(sum(pos) + sum(zero) + sum(neg) - 2 * min_)
