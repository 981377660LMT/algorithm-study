# 1190. 反转每对括号间的子串
# https://leetcode.cn/problems/reverse-substrings-between-each-pair-of-parentheses/solutions/795515/fan-zhuan-mei-dui-gua-hao-jian-de-zi-chu-gwpv/
# 按照从括号内到外的顺序，逐层反转每对匹配括号中的字符串，并返回最终的结果。
# O(n) 预处理翻转对应括号的位置
class Solution:
    def reverseParentheses(self, s: str) -> str:
        n = len(s)
        swapPair = [0] * n
        stack = []
        for i, v in enumerate(s):
            if v == "(":
                stack.append(i)
            elif v == ")":
                j = stack.pop()
                swapPair[i] = j
                swapPair[j] = i

        sb = []
        index, direction = 0, 1
        while index < n:
            if s[index] == "(" or s[index] == ")":
                index = swapPair[index]
                direction = -direction
            else:
                sb.append(s[index])
            index += direction
        return "".join(sb)


if __name__ == "__main__":
    print(Solution().reverseParentheses(s="a(bcdefghijkl(mno)p)q") == "apmnolkjihgfedcbq")
    print(Solution().reverseParentheses(s="((A)y)x"))
