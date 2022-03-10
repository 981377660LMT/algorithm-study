# 如果一个数 x 的约数之和 y（不包括他本身）比他本身小，那么 x 可以变成 y，y 也可以变成 x。
# 例如，4 可以变为 3，1 可以变为 7。
# 限定所有数字变换在不超过 n 的正整数范围内进行，求不断进行数字变换且不出现重复数字的最多变换步数。
# 1≤n≤50000

'''
树形DP, 树的最长路径的应用, 求一个森林里面所有树的最长路径的最大值
'''


from collections import defaultdict


def dfs(cur: int, parent: int) -> int:
    """求树的最长路径"""
    global res
    max1, max2 = 0, 0
    for next in adjMap[cur]:
        if next == parent:
            continue
        maxCand = dfs(next, cur) + 1
        if maxCand > max1:
            max1, max2 = maxCand, max1
        elif maxCand > max2:
            max2 = maxCand
    res = max(res, max1 + max2)
    return max1


n = int(input())
factorSums = [0] * (n + 1)
# 1. 因数筛求1-n中每个数的约数有哪些/约数之和
for factor in range(1, n + 1):
    # 约数不能是自己，所以从2开始枚举，如果约数可以是自己，从1开始枚举
    for multi in range(2, n // factor + 1):
        factorSums[factor * multi] += factor

# 2. 建图，并标记树根
notRoot = [False] * (n + 1)
adjMap = defaultdict(set)
# 注意从2开始，1的约数之和为0，不能要
for i in range(2, n + 1):
    if i > factorSums[i]:
        adjMap[factorSums[i]].add(i)
        notRoot[i] = True

# 3. 从森林的各个树的根节点dfs
res = 0
for i in range(1, n + 1):
    if notRoot[i]:
        continue
    dfs(i, -1)
print(res)
