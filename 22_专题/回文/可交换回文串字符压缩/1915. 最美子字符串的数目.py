# 如果某个字符串中 至多一个 字母出现 奇数 次，则称其为 最美 字符串
# 请你返回 word 中 最美非空子字符串 的数目
# 如果同样的子字符串在 word 中出现多次，那么应当对 每次出现 分别计数。

# 总结：
# 只考虑每个字母频率的奇偶性，就只考虑模2后的结果即可=>01状态压缩
# mask & 1 means whether it has odd 'a'
# mask & 2 means whether it has odd 'b'
# mask & 4 means whether it has odd 'c'
# word 由从 'a' 到 'j' 的小写英文字母组成
# 1 <= word.length <= 105


class Solution:
    def wonderfulSubstrings(self, word: str) -> int:
        preSum = [0 for _ in range(1 << 10)]
        preSum[0] = 1

        res = 0
        curState = 0
        for char in word:
            curState ^= 1 << (ord(char) - ord('a'))
            # 全为偶数的情况(异或全部抵消)
            res += preSum[curState]
            # 只有一个数出现奇数次的情况(仅一个位异或为1)
            for i in range(10):
                preState = curState ^ (1 << i)
                res += preSum[preState]

            preSum[curState] += 1

        return res


print(Solution().wonderfulSubstrings(word="aba"))
# 输出：4
# 解释：4 个最美子字符串如下所示：
# - "aba" -> "a"
# - "aba" -> "b"
# - "aba" -> "a"
# - "aba" -> "aba"

