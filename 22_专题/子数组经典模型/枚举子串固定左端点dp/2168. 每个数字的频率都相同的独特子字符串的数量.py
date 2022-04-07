class Solution:
    def equalDigitFrequency(self, s: str) -> int:
        n = len(s)
        visited = set()

        for i in range(n):
            counter = [0] * 10
            for j in range(i, n):
                counter[ord(s[j]) - ord('0')] += 1
                if len(set(counter) - {0}) == 1:
                    visited.add(s[i : j + 1])

        return len(visited)


print(Solution().equalDigitFrequency(s="12321"))
