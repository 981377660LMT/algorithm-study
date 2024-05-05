# 数组染色 使得任意两个数和为偶数
# 求染色方案数模10^9+7


# MOD = int(1e9 + 7)


# n = int(input())
# nums = list(map(int, input().split()))
# odd, even = [], []
# for i in nums:
#     if i % 2 == 0:
#         even.append(i)
#     else:
#         odd.append(i)
# n, m = len(odd), len(even)
# print((pow(2, n, MOD) + pow(2, m, MOD) - n - m - 2) % MOD)

mapping = {
    "r": 1,
    "e": 2,
    "d": 3,
}

s = input()
stack = []
res = 0
for i, c in enumerate(s):
    num = mapping[c]
    if stack and stack[-1] == num:
        stack.pop()
        res += 1
        if i + 1 < len(s):
            next = mapping[s[i + 1]]
            stack.append(num ^ next)
    else:
        stack.append(num)
print(res)
