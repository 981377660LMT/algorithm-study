# 字符串可以用 缩写 进行表示，缩写 的方法是将任意数量的 不相邻 的子字符串替换为相应子串的长度
# 0. 记录数字长度和字符下标
# 1. 数字没有'00'这种情况
# 2. 字符下标加上数字长度必须要对应 且不超出范围
# 3. 字符下标加数字长度最后等于word长度


class Solution:
    def validWordAbbreviation(self, word: str, abbr: str) -> bool:
        """
        双指针线性扫描：
        i 遍历 word，j 遍历 abbr。
        遇到字母，逐一比较；
        遇到数字，先检查不能是 '0' 开头，再累积成整数跳过相应数量的字符。
        最后要求两指针同时到达各自字符串末尾。
        时间 O(n + m)，空间 O(1)。
        """
        i, j = 0, 0
        n, m = len(word), len(abbr)

        while i < n and j < m:
            c = abbr[j]
            if c.isalpha():
                if word[i] != c:
                    return False
                i += 1
                j += 1
            else:
                if c == "0":
                    return False
                num = 0
                while j < m and abbr[j].isdigit():
                    num = num * 10 + (ord(abbr[j]) - ord("0"))
                    j += 1
                i += num

        return i == n and j == m


if __name__ == "__main__":
    sol = Solution()
    tests = [
        ("internationalization", "i12iz4n", True),
        ("apple", "a2e", False),
        ("substitution", "s10n", True),
        ("substitution", "s010n", False),
        ("substitution", "s0ubstitution", False),
    ]
    for w, a, exp in tests:
        res = sol.validWordAbbreviation(w, a)
        print(f"{w!r}, {a!r} -> {res} (expected {exp})")
