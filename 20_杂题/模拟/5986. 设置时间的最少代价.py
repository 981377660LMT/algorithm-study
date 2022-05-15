class Solution:
    def minCostSetTime(self, startAt: int, moveCost: int, pushCost: int, targetSeconds: int) -> int:
        """请你能返回设置 targetSeconds 秒钟加热时间需要花费的最少代价。
        
        将手指移到 任何其他数字 ，需要花费 moveCost 的单位代价。
        每 输入你手指所在位置的数字一次，需要花费 pushCost 的单位代价。
        """

        def cal(mins: int, secs: int) -> int:
            if 0 <= secs <= 99 and 0 <= mins <= 99:
                target = str(mins * 100 + secs)  # 注意这里的转换
                cur = str(startAt)

                res = 0
                for char in target:
                    if cur == char:
                        res += pushCost
                    else:
                        res += pushCost + moveCost
                        cur = char
                return res

            return int(1e20)

        mins, secs = divmod(targetSeconds, 60)
        # 微波炉显示的时间
        return min(cal(mins, secs), cal(mins - 1, secs + 60))


print(Solution().minCostSetTime(startAt=1, moveCost=2, pushCost=1, targetSeconds=600))
print(Solution().minCostSetTime(startAt=0, moveCost=1, pushCost=2, targetSeconds=76))
