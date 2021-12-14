from collections import Counter


class Solution:
    def numKLenSubstrNoRepeats(self, s: str, k: int) -> int:
        res = 0
        counter = Counter()
        overlap = 0
        for i in range(len(s)):
            if counter[s[i]] == 1:
                overlap += 1
            counter[s[i]] += 1

            if i - k >= 0:
                if counter[s[i - k]] == 2:
                    overlap -= 1
                counter[s[i - k]] -= 1

            if i >= k - 1:
                if overlap == 0:
                    res += 1

        return res


print(Solution().numKLenSubstrNoRepeats(s="havefunonleetcode", k=5))
