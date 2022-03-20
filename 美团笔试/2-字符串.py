# 小美正在点外卖。外卖平台上有两种优惠机制：折扣和满减。
# 用户在下单的时候可以选择这两种机制中的一种。
# 折扣机制即为部分商品提供了低于原价的优惠价，在结算的时候按优惠价计价。
# 满减机制则由一些规则组成，每条规则都形如“满c元减d元”，
# 意为若客户结算时所有商品的原价之和不低于c元，则客户可以减免d元。
# 用户只能使用一条满减规则，即若客户结算的商品总价为e元，
# 使用了满c元减d元的规则且e不小于c，则用户只需要支付e-d元。

# 小美有n种备选商品，每种商品均有原价和折扣价。现在小美想知道，
# 若仅购买前1,2,3….直到n种备选商品，是满减机制划算还是折扣机制划算。


from math import ceil


while True:
    try:
        n, optType = map(int, input().split())
        string = input()

        # n, optType = [6, 1]
        # string = 'hhhaaa'

        # n, optType = [6, 2]
        # string = 'hahaha'

        if optType == 1:
            res = []
            while len(res) < n:
                index = ceil(len(string) / 2) - 1
                res.append(string[index])
                string = string[:index] + string[index + 1 :]
            print(''.join(res))
        elif optType == 2:
            res = ''
            chars = list(string)
            while len(res) < n:
                char = chars.pop()
                index = len(res) // 2
                res = res[:index] + char + res[index:]
            print(res)
    except EOFError:
        break
