# () 得 1 分。 (()基础分为1)
# AB 得 A + B 分，其中 A 和 B 是平衡括号字符串。(可加性)
# (A) 得 2 * A 分，其中 A 是平衡括号字符串。(即level+1 之后 ()分数翻倍)

# 总结:处理level 碰到"()"结算
class Solution:
    def scoreOfParentheses1(self, s: str) -> int:
        return eval(s.replace(")(", ")+(").replace("()", "1").replace("(", "2*("))

    def scoreOfParentheses(self, s: str) -> int:
        n, res = len(s), 0
        level = 0
        for i in range(n):
            level += 1 if s[i] == "(" else -1
            if s[i] == ")" and i >= 1 and s[i - 1] == "(":
                res += 2**level

        return res


print(Solution().scoreOfParentheses("(())"))
print(Solution().scoreOfParentheses("(()(()))"))
