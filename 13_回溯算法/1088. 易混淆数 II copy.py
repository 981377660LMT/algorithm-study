# 给出正整数 N，请你返回 1 到 N 之间易混淆数字的数量。
# 1 <= N <= 10^9
pair = dict(zip('01689', '01986'))


class Solution:
    def confusingNumberII(self, n: int) -> int:
        def check(s: str) -> bool:
            return s[0] != '0' and ''.join([pair[char] for char in s])[::-1] != s

        def dfs(cur: int) -> None:
            nonlocal res
            if cur > n:
                return
            if check(str(cur)):
                res += 1
            for next in (0, 1, 6, 8, 9):
                dfs(cur * 10 + next)

        res = 0
        for start in (1, 6, 8, 9):
            dfs(start)
        return res


print(Solution().confusingNumberII(20))
