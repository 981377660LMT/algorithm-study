# 注意：所有小写字母的 ASCII 码的二进制第六位（从右向左）都是 1 ，而对应的大写字母第六位为 0 ，其他位都一样。

# !大写变小写、小写变大写 : 字符 ^= 32;
# !大写、小写变小写 : 字符 |= 32;
# 65 | 32 转为二进制（按8位来算）可以得到 0100 0001 | 0010 0000 = 0110 0001 = 97 = a


# !小写、大写变大写 : 字符 &= -33;
# 转换为大写即与0b1...011111（-33）做与运算


class Solution:
    def toLowerCase(self, s: str) -> str:
        return "".join([char if not char.isalpha() else chr((ord(char) | 32)) for char in s])


print(Solution().toLowerCase("Hello"))
