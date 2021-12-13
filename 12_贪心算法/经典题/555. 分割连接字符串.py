from typing import List

# 顺序连接+可内部翻转+最后可在循环字符串的某个位置分割它;
# 找到字典序最大的字符串。
class Solution:
    def splitLoopedString(self, strs: List[str]) -> str:
        # maxChar必须放到首位
        maxChar = max(max(s) for s in strs)
        # 每个字符串必须最大
        strs = [max(s, s[::-1]) for s in strs]
        res = ''
        # 枚举删除strs[i]
        for i, word in enumerate(strs):
            other = ''.join(strs[i + 1 :] + strs[:i])
            # 枚举分割位置j
            for j in range(len(word)):
                if word[j] == maxChar:
                    res = max(
                        res,
                        word[j:] + other + word[:j],
                        word[: j + 1][::-1] + other + word[j + 1 :][::-1],
                    )

        return res


print(Solution().splitLoopedString(["abc", "xyz"]))
# 输出: "zyxcba"
# 解释: 你可以得到循环字符串 "-abcxyz-", "-abczyx-", "-cbaxyz-", "-cbazyx-"，
# 其中 '-' 代表循环状态。
# 答案字符串来自第四个循环字符串，
# 你可以从中间字符 'a' 分割开然后得到 "zyxcba"。

