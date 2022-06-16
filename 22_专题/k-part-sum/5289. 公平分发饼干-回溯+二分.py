from typing import List

# !二分确定范围之后 便于在回溯时剪枝


class Solution:
    def distributeCookies(self, cookies: List[int], k: int) -> int:
        def check(mid: int) -> bool:
            """分割成k组,各组和的最大值小于等于mid"""

            def dfs(index: int, groups: List[int]) -> bool:
                if index == n:
                    return True
                for i in range(k):
                    if groups[i] + cookies[index] <= mid:
                        groups[i] += cookies[index]
                        if dfs(index + 1, groups):
                            return True
                        groups[i] -= cookies[index]
                return False

            return dfs(0, [0] * k)

        cookies.sort(reverse=True)
        n = len(cookies)
        left, right = max(cookies), sum(cookies)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left
