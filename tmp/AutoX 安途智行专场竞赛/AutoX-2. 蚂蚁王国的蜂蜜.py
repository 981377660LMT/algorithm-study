from collections import defaultdict
from typing import List

# 若 type 为 1，表示获取了一份价格为 value 的报价
# 若 type 为 2，表示删除了一个价格为 value 的报价
# 若 type 为 3，表示计算当前蜂蜜的均价；若当前不存在任何有效报价，返回 -1
# 若 type 为 4，表示计算当前价格的方差；若当前不存在任何有效报价，返回 -1


class Solution:
    def honeyQuotes(self, handle: List[List[int]]) -> List[float]:
        counter, count, sum_ = defaultdict(int), 0, 0
        res = []
        for type, *rest in handle:
            if type == 1:
                value = rest[0]
                counter[value] += 1
                count += 1
                sum_ += value
            elif type == 2:
                value = rest[0]
                counter[value] -= 1
                count -= 1
                sum_ -= value
            elif type == 3:
                if count == 0:
                    res.append(-1)
                else:
                    res.append(sum_ / count)
            elif type == 4:
                if count == 0:
                    res.append(-1)
                else:
                    avg = sum_ / count
                    diff = 0
                    for k in counter:
                        diff += (k - avg) ** 2 * counter[k]
                    res.append(diff / count)

        return res


print(Solution().honeyQuotes(handle=[[3], [1, 10], [1, 0], [3], [4], [2, 10], [3]]))
