from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)


# 当你对上对手时，需要在经验和精力上都 严格 超过对手才能击败他们，然后在可能的情况下继续对上下一个对手。
# 击败第 i 个对手会使你的经验 增加 experience[i]，但会将你的精力 减少  energy[i] 。


class Solution:
    def minNumberOfHours(
        self, initialEnergy: int, initialExperience: int, energy: List[int], experience: List[int]
    ) -> int:
        """模拟 如果遇到瓶颈就加delta"""
        res = 0
        for a, b in zip(energy, experience):
            if initialEnergy <= a:
                delta = a + 1 - initialEnergy
                res += delta
                initialEnergy = a + 1
            if initialExperience <= b:
                delta = b + 1 - initialExperience
                res += delta
                initialExperience = b + 1
            initialEnergy -= a
            initialExperience += b
        return res

    def minNumberOfHours2(
        self, initialEnergy: int, initialExperience: int, energy: List[int], experience: List[int]
    ) -> int:
        sum1 = sum(energy)
        res1 = max(0, sum1 + 1 - initialEnergy)

        # 需要多少经验
        def check(mid: int) -> bool:
            cur = mid
            for exp in experience:
                if cur <= exp:
                    return False
                cur += exp
            return True

        left, right = 0, int(1e9)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return max(0, left - initialExperience) + res1


print(
    Solution().minNumberOfHours(
        initialEnergy=5, initialExperience=3, energy=[1, 4, 3, 2], experience=[2, 6, 3, 1]
    )
)
