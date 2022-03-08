# 对于每一根钢筋 x，我们会写下 +x，-x 或者 0。
# 0 <= rods.length <= 20
# 暴力枚举3^20肯定会TLE 需要折半枚举
from collections import defaultdict
from typing import DefaultDict, Dict, List, Set, Tuple

# 你正在安装一个广告牌，并希望它高度最大。这块广告牌将有两个钢制支架，两边各一个。每个钢支架的高度必须相等。
# 你有一堆可以焊接在一起的钢筋 rods。举个例子，如果钢筋的长度为 1、2 和 3，则可以将它们焊接在一起形成长度为 6 的支架。
# 返回 广告牌的最大可能安装高度 。如果没法安装广告牌，请返回 0 。
# 题意可以等价于在列表中选出一个数字集合，集合中的每个数字可以乘上1或者-1，乘完之后集合中的所有数字之和为0，求满足这个条件的集合正数之和的最大值
# 枚举数字集合使得和为0是个经典问题，可以用折半搜索来做，时间复杂度为(3^(n/2))
# 时间复杂度：O(3^N/2)


class Solution:
    def tallestBillboard(self, rods: List[int]) -> int:
        # 枚举n/2*3^n/2超时
        # def getSum(nums: List[int]) -> Dict[int, int]:
        #     n = len(nums)
        #     res = {}
        #     for state in range(3 ** n):
        #         curSum, absSum = 0, 0
        #         for i in range(n):
        #             mod = (state // (3 ** i)) % 3
        #             if mod == 1:
        #                 curSum += nums[i]
        #                 absSum += nums[i]
        #             elif mod == 2:
        #                 curSum -= nums[i]
        #                 absSum += nums[i]
        #         res[curSum] = max(res.get(curSum, 0), absSum)
        #     return res

        # dfs 3^(n/2) 不超时 还是dfs(index,curSum) 比较好
        def getSum(nums: List[int]) -> Dict[int, int]:
            def dfs(index: int, curSum: int, absSum: int) -> None:
                if index == n:
                    res[curSum] = max(res.get(curSum, 0), absSum)
                    return
                dfs(index + 1, curSum, absSum)
                dfs(index + 1, curSum - nums[index], absSum + nums[index])
                dfs(index + 1, curSum + nums[index], absSum + nums[index])

            n = len(nums)
            res = {}
            dfs(0, 0, 0)
            return res

        half = len(rods) // 2
        d1, d2 = getSum(rods[:half]), getSum(rods[half:])
        if len(d1) > len(d2):
            d1, d2 = d2, d1

        res = 0
        for curSum, absSum in d1.items():
            if curSum in d2:
                res = max(res, (absSum + d2[curSum]) // 2)
        return res


print(Solution().tallestBillboard([1, 2, 3, 6]))
print(Solution().tallestBillboard([1, 2, 3, 4, 5, 6]))
