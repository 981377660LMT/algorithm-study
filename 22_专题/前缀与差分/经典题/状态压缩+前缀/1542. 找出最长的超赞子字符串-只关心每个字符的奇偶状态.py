# 请返回 s 中`最长`的 超赞子字符串 的长度。
# !即可交换变成回文：奇数的频率的字母种数不能超过1


# 因为我们只关心奇偶，我们可以用一个长度为 10 的 0-1 序列来表示任意一个子串 s'
# s 仅由数字组成 => 状态压缩
# 1 <= s.length <= 10^5

from collections import defaultdict


# 1915. 最美子字符串的数目
class Solution:
    def longestAwesome(self, s: str) -> int:
        # 每个状态最早出现的位置
        first = defaultdict(lambda: int(1e20))
        first[0] = -1

        res = 0
        curState = 0
        for i, char in enumerate(s):
            curState ^= 1 << (ord(char) - ord('0'))
            res = max(res, i - first[curState])
            for diff in range(10):
                pre = curState ^ 1 << (diff)
                res = max(res, i - first[pre])

            first[curState] = min(first[curState], i)

        return res


print(Solution().longestAwesome("3242415"))
# 输入：s = "3242415"
# 输出：5
# 解释："24241" 是最长的超赞子字符串，交换其中的字符后，可以得到回文 "24142"
