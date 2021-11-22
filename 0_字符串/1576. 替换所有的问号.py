from string import ascii_lowercase

# 给你一个仅包含小写英文字母和 '?' 字符的字符串 s，
# 请你将所有的 '?' 转换为若干小写字母，
# 使最终的字符串不包含任何 连续重复 的字符。
# 理论上三个字母就可以替换所有


class Solution:
    def modifyString(self, s: str) -> str:
        # sentinel
        sb = list('#' + s + '#')

        for index in range(len(sb)):
            if sb[index] != '?':
                continue

            for choose in ascii_lowercase:
                if sb[index - 1] != choose and sb[index + 1] != choose:
                    sb[index] = choose
                    break

        return ''.join(sb[1:-1])


# 输入：s = "j?qg??b"
# 输出："jaqgacb"
