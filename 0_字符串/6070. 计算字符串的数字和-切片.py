MOD = int(1e9 + 7)
INF = int(1e20)

# 将 s 拆分 成长度为 k 的若干 连续数字组 ，
# 使得前 k 个字符都分在第一组，接下来的 k 个字符都分在第二组，
# 依此类推。注意，最后一个数字组的长度可以小于 k 。


class Solution:
    def digitSum(self, s: str, k: int) -> str:
        while len(s) > k:
            groups = [s[i : i + k] for i in range(0, len(s), k)]
            chars = []
            for g in groups:
                sum_ = sum(int(c) for c in g)
                chars.append(str(sum_))
            s = ''.join(chars)
        return s

