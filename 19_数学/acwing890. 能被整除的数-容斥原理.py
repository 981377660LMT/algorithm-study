# 给定一个整数 n 和 m 个不同的质数 p1,p2,…,pm。
# 请你求出 1∼n 中能被 p1,p2,…,pm 中的至少一个数整除的整数有多少个。

'''
容斥原理求解方案数
'''

n, m = map(int, input().split())
arr = list(map(int, input().split()))

# 二进制枚举所有集合的选择情况
res = 0
for state in range(1, (1 << m)):
    multi = 1
    for i in range(m):
        if state & (1 << i):
            multi *= arr[i]

    # 因为都是质数，所以乘积就是最小公倍数，n//最小公倍数 就是集合包含能同时被这些质数整除的数值数量
    if bin(state)[2:].count('1') & 1:
        res += n // multi
    else:
        res -= n // multi

print(res)

