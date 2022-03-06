# 奶牛们正在试图找到自己在这个堆叠中应该所处的位置顺序。

# 一头牛支撑不住的可能性取决于它头上所有牛的总重量（不包括它自己）减去它的身体强壮程度的值，
# 现在称该数值为风险值，风险值越大，
# 这只牛撑不住的可能性越高。

# 您的任务是确定奶牛的排序，使得所有奶牛的风险值中的最大值尽可能的小。

# 输出一个整数，表示最大风险值的最小可能值。


# 发现一个指标:
# Wi + Si
# 这个指标越小, 就应该放在前面, 越大, 就应该放在后面, 从而使最大风险减少.

n = int(input())
cows = []
for i in range(n):
    weight, strong = map(int, input().split())
    cows.append((weight + strong, weight, strong))
cows.sort()
res = -int(1e20)
curSum = 0
for _, weight, strong in cows:
    res = max(res, curSum - strong)
    curSum += weight
print(res)
