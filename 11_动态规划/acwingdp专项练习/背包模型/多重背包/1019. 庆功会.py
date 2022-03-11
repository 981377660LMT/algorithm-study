# 为了庆贺班级在校运动会上取得全校第一名成绩，班主任决定开一场庆功会，为此拨款购买奖品犒劳运动员。
# 期望拨款金额能购买最大价值的奖品，可以补充他们的精力和体力。


# 第一行二个数n，m，其中n代表希望购买的奖品的种数，m表示拨款金额。
# 接下来n行，每行3个数，v、w、s，分别表示第I种奖品的价格、价值（价格与价值是不同的概念）和能购买的最大数量（买0件到s件均可）。
# 输出格式
# 一行：一个数，表示此次购买能获得的最大的价值（注意！不是价格）。

n, cap = map(int, input().split())
goods = []
for i in range(n):
    cost, score, count = map(int, input().split())
    cur = 1
    while cur <= count:
        goods.append([cost * cur, score * cur])
        count -= cur
        cur <<= 1
    if count:
        goods.append([cost * count, score * count])

dp = [0] * (cap + 1)
for cost, score in goods:
    for i in range(cap, cost - 1, -1):
        dp[i] = max(dp[i], dp[i - cost] + score)
print(dp[-1])
