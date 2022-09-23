# !数字的各位组成两个整数乘积的最大值
# 保证输入的数字至少包含两个非0数字
# O(LlogL) 贪心，分奇偶,两个数的差要尽量小

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    digits = sorted(input(), reverse=True)
    A, B = digits[::2], digits[1::2]
    for i in range(min(len(A), len(B))):
        if A[i] != B[i]:  # 交换第一个不相等的位置 两个数的差要尽量小
            A[i], B[i] = B[i], A[i]
            break
    print(int("".join(A)) * int("".join(B)))
