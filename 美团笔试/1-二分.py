# 小美正在点外卖。外卖平台上有两种优惠机制：折扣和满减。
# 用户在下单的时候可以选择这两种机制中的一种。
# 折扣机制即为部分商品提供了低于原价的优惠价，在结算的时候按优惠价计价。
# 满减机制则由一些规则组成，每条规则都形如“满c元减d元”，
# 意为若客户结算时所有商品的原价之和不低于c元，则客户可以减免d元。
# 用户只能使用一条满减规则，即若客户结算的商品总价为e元，
# 使用了满c元减d元的规则且e不小于c，则用户只需要支付e-d元。

# 小美有n种备选商品，每种商品均有原价和折扣价。现在小美想知道，
# 若仅购买前1,2,3….直到n种备选商品，是满减机制划算还是折扣机制划算。
from bisect import bisect_right
from itertools import accumulate


while True:
    try:
        n = int(input())
        prices = list(map(int, input().split()))
        discounts = list(map(int, input().split()))

        # 代表满减规则的数量。
        m = int(input())
        threshold = list(map(int, input().split()))
        minus = list(map(int, input().split()))

        # n = 3
        # prices = [5, 10, 8]
        # discounts = [5, 8, 7]
        # m = 2
        # threshold = [15, 22]
        # minus = [1, 4]

        # 输出一个长度为 n 的字符串，若仅购买前 i 种备选商品时满减机制更划算则第 i 个字符为M，若折扣机制更划算则第 i 个字符为Z，若两种机制带来的优惠相同则输出B。
        res = []
        plan1 = list(accumulate(discounts))
        plan2 = []
        # 用户只能使用一条满减规则
        curSum = 0
        curMinus = 0
        for i in range(n):
            curSum += prices[i]
            index = bisect_right(threshold, curSum) - 1
            curMinus = 0 if index < 0 else minus[index]
            plan2.append(curSum - curMinus)
        for n1, n2 in zip(plan1, plan2):
            if n1 > n2:
                res.append('M')
            elif n1 == n2:
                res.append('B')
            else:
                res.append('Z')

        print(''.join(res))

    except EOFError:
        break
