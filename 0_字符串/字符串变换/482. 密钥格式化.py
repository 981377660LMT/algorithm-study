# 482. 密钥格式化
# https://leetcode.cn/problems/license-key-formatting/description/
#
# 给定一个许可密钥字符串 s，仅由字母、数字字符和破折号组成。
# 字符串由 n 个破折号分成 n + 1 组。你也会得到一个整数 k 。
# 我们想要重新格式化字符串 s，使每一组包含 k 个字符，除了第一组，它可以比 k 短，但仍然必须包含至少一个字符。
# 此外，两组之间必须插入破折号，并且应该将所有小写字母转换为大写字母。
# 返回 重新格式化的许可密钥 。


class Solution:
    def licenseKeyFormatting(self, s: str, k: int) -> str:
        chars = [c.upper() for c in s if c != "-"]
        n = len(chars)
        groups = []
        for end in range(n, 0, -k):
            start = max(0, end - k)
            groups.append("".join(chars[start:end]))
        return "-".join(reversed(groups))
