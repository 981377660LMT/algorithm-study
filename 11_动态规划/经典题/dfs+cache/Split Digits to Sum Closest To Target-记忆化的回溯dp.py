# 将数字字符串任意拆分,最小化他们的和与target的距离

# 带记忆化的回溯dp 好处是可以剪枝
from functools import lru_cache


class Solution:
    def solve(self, s: str, target: int):
        @lru_cache(None)
        def dfs(index: int, curSum: int) -> None:
            if index == n:
                self.res = min(self.res, abs(curSum - target))
                return

            # 剪枝
            if curSum >= target and curSum - target > self.res:
                return

            for len_ in range(4):
                if index + len_ >= n:
                    break
                dfs(index + len_ + 1, curSum + int(s[index : index + len_ + 1]))

        n = len(s)
        nums = [int(char) for char in s]
        sum_ = sum(nums)
        if sum_ >= target:
            return sum_ - target

        self.res = abs(sum_ - target)
        dfs(0, 0)
        return self.res


print(Solution().solve(s="112", target=10))
# We can partition s into "1" + "12" which sums to 13 and abs(13 - 10) = 3.


# 1 ≤ len(s), target ≤ 1,000
class Solution2:
    def solve(self, s, target):
        @lru_cache(None)
        def dfs(index: int, diff: int) -> int:
            if index >= n:
                return abs(diff)
            if nums[index] == 0:
                return dfs(nextNoneZero[index], diff)

            res = cand
            cur = 0
            for i in range(index, n):
                cur = cur * 10 + nums[i]
                res = min(res, dfs(i + 1, diff - cur))
            return res

        n = len(s)
        nums = [int(char) for char in s]

        nextNoneZero = [int(1e20)] * n
        for i in range(n - 1, -1, -1):
            if nums[i] != 0:
                nextNoneZero[i] = i
            elif i + 1 < n:
                nextNoneZero[i] = nextNoneZero[i + 1]

        cand = abs(target - sum(nums))
        return dfs(0, target)
