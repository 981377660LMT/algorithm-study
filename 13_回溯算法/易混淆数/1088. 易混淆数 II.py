# 给出正整数 N，请你返回 1 到 N 之间易混淆数字的数量。
# 1 <= N <= 10^9
class Solution:
    def confusingNumberII(self, n: int) -> int:
        def check(s: str):
            return s[0] != '0' and ''.join([pair[char] for char in s])[::-1] != s

        pair = dict(zip('01689', '01986'))
        res, arr = 0, ['']

        for _ in range(len(str(n))):
            arr = [pre + cur for pre in arr for cur in '01689' if int(pre + cur) <= n]
            res += sum(check(s) for s in arr)

        return res


print(Solution().confusingNumberII(20))
