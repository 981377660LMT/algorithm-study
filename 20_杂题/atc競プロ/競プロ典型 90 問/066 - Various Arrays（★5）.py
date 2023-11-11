import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

# !每个nums[i] 都是介于 [Li,Ri]的一个randint
# 求逆序对的期望值
# 1≤N≤100
# L,R<=100
# 枚举每个comb 计算贡献 O(n^2*max(L,R))

# !期望值的线性可叠加性 (类似计算贡献)

n = int(input())
lefts = []
rights = []
for _ in range(n):
    left, right = map(int, input().split())
    lefts.append(left)
    rights.append(right)

res = 0
for i in range(n):
    len1 = rights[i] - lefts[i] + 1
    for j in range(i + 1, n):
        len2 = rights[j] - lefts[j] + 1
        score = 0  # 贡献
        for k in range(lefts[i], rights[i] + 1):
            score += max(0, min(len2, k - lefts[j]))
        res += score / (len1 * len2)

print(res)
