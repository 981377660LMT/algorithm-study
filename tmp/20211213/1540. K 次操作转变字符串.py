# 你的目标是在 k 次操作以内把字符串 s 转变成 t 。
class Solution:
    def canConvertString(self, s: str, t: str, k: int) -> bool:
        if len(s) != len(t):
            return False

        n = len(s)
        dist = [0 for _ in range(26)]  # 这个距离需要操作几次
        for i in range(n):
            diff = min((ord(t[i]) - ord(s[i])) % 26, (ord(s[i]) - ord(t[i]) % 26))
            # 第diff轮需要的总操作次数
            dist[diff] += 1
        print(dist)

        for round in range(1, 26):
            # 如果其中的一个下标在第 i 次操作时进行了切换，另一个下标必须等到第 i+26 次操作时才能进行切换
            maxDist = round + 26 * (dist[round] - 1)  # 有重复不要紧，下一轮，下下轮，下下下轮
            if maxDist > k:  # 超了， False
                return False

        return True


print(Solution().canConvertString(s="input", t="ouput", k=9))
# 输出：true
# 解释：第 6 次操作时，我们将 'i' 切换 6 次得到 'o' 。
# 第 7 次操作时，我们将 'n' 切换 7 次得到 'u' 。
