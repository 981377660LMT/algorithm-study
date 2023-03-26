# 给你一个整数 money ，表示你总共有的钱数（单位为美元）和另一个整数 children ，表示你要将钱分配给多少个儿童。

# 你需要按照如下规则分配：

# 所有的钱都必须被分配。
# 每个儿童至少获得 1 美元。
# 没有人获得 4 美元。
# 请你按照上述规则分配金钱，并返回 最多 有多少个儿童获得 恰好 8 美元。如果没有任何分配方案，返回 -1 。

# 1 <= money <= 200
# 2 <= children <= 30


class Solution:
    def distMoney(self, money: int, children: int) -> int:
        if money < children:
            return -1

        # !枚举8美元的儿童数
        res = -1
        for i in range(children + 1):
            restMoney = money - i * 8
            restChildren = children - i
            if restMoney < restChildren:
                break
            if restMoney == 4 and restChildren == 1:
                break
            if restMoney > 0 and restChildren == 0:
                break
            res = max(res, i)  # 剩余的钱分给restChildren个儿童

        return res
