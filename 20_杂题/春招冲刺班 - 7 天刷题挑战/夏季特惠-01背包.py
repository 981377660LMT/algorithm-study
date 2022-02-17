# 1 <= n <= 500
# 0 <= x <= 10,000

# 某公司游戏平台的夏季特惠开始了，
# 你决定入手一些游戏。现在你一共有X元的预算，
# 该平台上所有的n个游戏均有折扣，标号为i的游戏的原价a;元，
# 现价只要 b元(也就是说该游戏可以优惠;-b元)并且你购买该游戏能获得快乐值为u;。
# 由于优惠的存在，你可能做出一些冲动消费导致最终买游戏的总费用超过预算，
# 但只要满足获得的总优惠金额不低于超过预算的总金额，那在心理上就不会觉得吃亏。
# 现在你希望在心理上不觉得吃亏的前提下，获得尽可能多的快乐值。

# 由于优惠的存在，你可能做出一些冲动消费导致最终买游戏的总费用超过预算，
# 但只要满足获得的总优惠金额不低于超过预算的总金额，那在心理上就不会觉得吃亏
# 这句话等价于售价为b-(a-b),注意价格会出现负数，要预先加上这些


n, capacity = map(int, input().split())
games = []

overflow = 0
for _ in range(n):
    a, b, score = map(int, input().split())
    price = b - (a - b)
    if price < 0:
        # 注意这里选择后背包容量相当于变大了
        capacity += -price
        overflow += score
    else:
        games.append((price, score))

res = 0
dp = [0] * (capacity + 1)
for price, score in games:
    for i in range(capacity, -1, -1):
        if i >= price:
            dp[i] = max(dp[i], dp[i - price] + score)
            res = max(res, dp[i])

print(res + overflow)
