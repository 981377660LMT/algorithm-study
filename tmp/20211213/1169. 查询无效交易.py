from typing import List
from collections import defaultdict

# 如果出现下述两种情况，交易 可能无效：

# 交易金额超过 ¥1000
# 或者，它和另一个城市中同名的另一笔交易相隔不超过 60 分钟（包含 60 分钟整）
# 名称，时间（以分钟计），金额以及城市。
# 返回可能无效的交易列表
class Solution:
    def invalidTransactions(self, transactions: List[str]) -> List[str]:
        trans = [x.split(',') for x in transactions]
        res = []

        # 每次循环只检验i是不是合法的
        for i in range(len(trans)):
            name, time, money, city = trans[i]
            time = int(time)
            if int(money) > 1000:
                res.append(transactions[i])
                continue

            for j in range(len(trans)):
                if i == j:
                    continue
                name1, time1, _, city1 = trans[j]
                if name1 == name and city1 != city and abs(int(time1) - time) <= 60:
                    res.append(transactions[i])
                    break

        return res


print(Solution().invalidTransactions(transactions=["alice,20,800,mtv", "alice,50,100,beijing"]))
# 输出：["alice,20,800,mtv","alice,50,100,beijing"]
# 解释：第一笔交易是无效的，
# 因为第二笔交易和它间隔不超过 60 分钟、
# 名称相同且发生在不同的城市。同样，第二笔交易也是无效的。
