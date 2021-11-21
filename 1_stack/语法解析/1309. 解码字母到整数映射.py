# 反向遍历


class Solution:
    def freqAlphabets(self, s: str) -> str:
        sb, cursor = [], len(s) - 1
        while cursor >= 0:
            if s[cursor] == '#':
                sb.append(chr(int(s[cursor - 2 : cursor]) + 96))
                cursor -= 3
            else:
                sb.append(chr(int(s[cursor]) + 96))
                cursor -= 1

        return ''.join(reversed(sb))


print(Solution().freqAlphabets("1326#"))

# 字符（'a' - 'i'）分别用（'1' - '9'）表示。
# 字符（'j' - 'z'）分别用（'10#' - '26#'）表示。

