from collections import Counter


class Solution:
    def checkInclusion(self, s1: str, s2: str) -> bool:
        if len(s1) > len(s2):
            return False

        m, n = len(s1), len(s2)
        left = 0
        needCount = m
        needCounter = Counter(s1)

        for right in range(n):
            if needCounter[s2[right]] > 0:
                needCount -= 1
            needCounter[s2[right]] -= 1
            right += 1

            if needCount == 0:
                return True

            if right - left >= m:
                if needCounter[s2[left]] >= 0:
                    needCount += 1
                needCounter[s2[left]] += 1
                left += 1

        return False

