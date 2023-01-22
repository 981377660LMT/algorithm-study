# 给你两个下标从 0 开始的 二元 字符串 s 和 target ，两个字符串的长度均为 n 。你可以对 s 执行下述操作 任意 次：

# 选择两个 不同 的下标 i 和 j ，其中 0 <= i, j < n 。
# 同时，将 s[i] 替换为 (s[i] OR s[j]) ，s[j] 替换为 (s[i] XOR s[j]) 。
# 例如，如果 s = "0110" ，你可以选择 i = 0 和 j = 2，然后同时将 s[0] 替换为 (s[0] OR s[2] = 0 OR 1 = 1)，并将 s[2] 替换为 (s[0] XOR s[2] = 0 XOR 1 = 1)，最终得到 s = "1110" 。

# 如果可以使 s 等于 target ，返回 true ，否则，返回 false 。


# !找规律：
# 0 0 => 0 0
# 0 1 => 1 1
# 1 0 => 1 1
# 1 1 => 1 0
# 如果存在1的话,最后都可以变为1000...000


class Solution:
    def makeStringsEqual(self, s: str, target: str) -> bool:
        # 有1的话可变成 10000000...
        if "1" not in s:
            return s == target
        return "1" in target

        return ("1" in s) == ("1" in target)  # `同时`存在1或者同时不存在1
