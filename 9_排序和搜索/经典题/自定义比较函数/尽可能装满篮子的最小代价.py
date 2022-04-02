# Fruit Basket Packing
# Each basket can only contain fruits of the same kind
# Each basket should be as full as possible
# Each basket should be as cheap as possible

# 贪心
# 处理每种水果
class Solution:
    def solve(self, fruits, k, capacity):
        goods = []

        for cost, size, count in fruits:
            while count > 0:
                # 每个篮子装几个
                take = min(capacity // size, count)
                if take == 0:
                    break

                # 可以装几个篮子
                fill = count // take
                # (waste capacity, basket cost, number of such baskets we can fill using this option)
                goods.append((capacity - take * size, take * cost, fill))
                count -= fill * take

        res = 0
        # 最小浪费，最小花费，最小篮子数
        for _, basketCost, fill in sorted(goods):
            realFill = min(k, fill)
            res += basketCost * realFill
            k -= realFill
            if k == 0:
                break

        return res


print(Solution().solve(fruits=[[4, 2, 3], [5, 3, 2], [1, 3, 2]], k=2, capacity=4))
# [cost, size, total]
