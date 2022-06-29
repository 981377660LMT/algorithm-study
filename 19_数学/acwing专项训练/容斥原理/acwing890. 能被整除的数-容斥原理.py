# 给定一个整数 n 和 m 个不同的质数 p1,p2,…,pm。
# 请你求出 1∼n 中能被 p1,p2,…,pm 中的至少一个数整除的整数有多少个。
# m<=10
'''
容斥原理求解方案数
'''

n, m = map(int, input().split())
arr = list(map(int, input().split()))


res = 0
for state in range(1, (1 << m)):  # 枚举被哪些数整除
    mul = 1
    for i in range(m):
        if state & (1 << i):
            mul *= arr[i]

    # 因为都是质数，所以乘积就是最小公倍数，n//最小公倍数 就是集合包含能同时被这些质数整除的数值数量
    # !奇数个元素系数为 1，偶数个元素为 -1
    if state.bit_count() & 1:
        res += n // mul
    else:
        res -= n // mul

print(res)

