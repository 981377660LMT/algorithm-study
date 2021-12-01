# 利用当前出现过的数字构造`下一个`距离当前时间最近的时刻

# 对四个位置的数暴力枚举所有的情况，如果组成的时间合法并且刚好大于time，就找到结果了
# 循环结束还没有找到，说明time是`最大的排列组合`，返回非0`最小值`的4个位置重复(明天最早时间)
# 1. 注意:sorted预处理很关键 `找到就直接是答案`
class Solution:
    def nextClosestTime(self, time: str) -> str:
        time = time[:2] + time[3:]
        nums = sorted(list(set(time)))

        for n1 in nums:
            for n2 in nums:
                for n3 in nums:
                    for n4 in nums:
                        candi = n1 + n2 + n3 + n4
                        if candi < '2400' and n3 < '6' and candi > time:
                            return n1 + n2 + ':' + n3 + n4

        c = nums[0] if nums[0] != 0 else nums[1]
        return c * 2 + ':' + c * 2


print(Solution().nextClosestTime("19:34"))
