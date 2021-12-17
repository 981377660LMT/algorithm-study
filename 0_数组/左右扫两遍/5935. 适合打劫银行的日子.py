from typing import List

# 如果第 i 天满足以下所有条件，我们称它为一个适合打劫银行的日子：
# 第 i 天前和后都分别至少有 time 天。
# 第 i 天前连续 time 天警卫数目都是非递增的。
# 第 i 天后连续 time 天警卫数目都是非递减的。

# 总结：扫两边，统计递增递减个数


class Solution:
    def goodDaysToRobBank(self, security: List[int], time: int) -> List[int]:
        n = len(security)

        pre = [0] * n
        for i in range(1, n):
            if security[i - 1] >= security[i]:
                pre[i] = pre[i - 1] + 1

        suffix = [0] * n
        for i in range(n - 2, -1, -1):
            if security[i] <= security[i + 1]:
                suffix[i] = suffix[i + 1] + 1

        res = []
        for i in range(time, n - time):
            if pre[i] >= time and suffix[i] >= time:
                res.append(i)

        return res


print(Solution().goodDaysToRobBank(security=[5, 3, 3, 3, 5, 6, 2], time=2))
# 输出：[2,3]
# 解释：
# 第 2 天，我们有 security[0] >= security[1] >= security[2] <= security[3] <= security[4] 。
# 第 3 天，我们有 security[1] >= security[2] >= security[3] <= security[4] <= security[5] 。
# 没有其他日子符合这个条件，所以日子 2 和 3 是适合打劫银行的日子
