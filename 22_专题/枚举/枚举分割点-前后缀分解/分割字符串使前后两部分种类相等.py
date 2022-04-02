class Solution:
    def solve(self, s: str) -> int:
        n = len(s)
        prefix = [0] * n
        suffix = [0] * n

        counter = set()
        for i in range(n):
            counter.add(s[i])
            prefix[i] = len(counter)

        counter = set()
        for i in range(n - 1, -1, -1):
            counter.add(s[i])
            suffix[i] = len(counter)

        return sum(prefix[i] == suffix[i + 1] for i in range(n - 1))


print(Solution().solve(s="abaab"))

# We can split it by "ab" + "aab" and "aba" + "ab"
