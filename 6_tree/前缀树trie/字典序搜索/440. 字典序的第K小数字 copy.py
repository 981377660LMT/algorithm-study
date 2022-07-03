class Solution:
    def findKthNumber(self, n: int, k: int) -> int:
        """求字典序第k小的数字,前序遍历"""

        def calStep(cur: int, target: int) -> int:
            """cur到target的距离，检查子树数量即可"""
            step = 0
            while cur <= n:
                step += min(target, n + 1) - cur
                cur *= 10
                target *= 10
            return step

        cur = 1
        k -= 1

        while k:
            # 向右走需要的步数
            steps = calStep(cur, cur + 1)
            # 向右走
            if steps <= k:
                k -= steps
                cur += 1
            # 向下走
            else:
                k -= 1
                cur *= 10

        return cur


print(Solution().findKthNumber(11, 2))
