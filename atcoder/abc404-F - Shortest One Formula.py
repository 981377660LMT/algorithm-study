# F - Shortest One Formula
# https://atcoder.jp/contests/abc403/editorial/12771
#
#
# 给定一个正整数 **N**，需要找到一个由以下字符组成的数式，使其值等于 **N**，并且作为字符串的长度最短。
#
# 数式只能由以下字符组成：
# - 数字 `1`
# - 运算符 `+` 和 `*`
# - 括号 `(` 和 `)`
#
# 需要满足以下条件：
# 1. 数式必须符合以下的 BNF 语法规则：
#    ```
#    <expr>   ::= <term> | <expr> "+" <term>
#    <term>   ::= <factor> | <term> "*" <factor>
#    <factor> ::= <number> | "(" <expr> ")"
#    <number> ::= "1" | "1" <number>
#    ```
# 2. 数式的值必须等于 **N**。
# 3. 数式作为字符串的长度必须最短。
#
# 例如：
# - 对于 **N = 4**，可能的数式有 `1+1+1+1` 或 `(1+1)*(1+1)`，但 `(1+1)*(1+1)` 更短，答案为 `(1+1)*(1+1)`。
# - 对于 **N = 6**，可能的数式有 `1+1+1+1+1+1` 或 `(1+1+1)*(1+1)`，答案为 `(1+1+1)*(1+1)`。
#
# **约束条件**：
# - 1 ≤ N ≤ 2000
#
# 输入：
# - 一个整数 **N**
#
# 输出：
# - 一个符合条件的最短数式字符串。


if __name__ == "__main__":

    def update(x: str, y: str) -> str:
        return x if len(x) < len(y) else y

    n = int(input())
    dp = ["1" * 1000] * (n + 1)
    dp_term = ["1" * 1000] * (n + 1)

    for i in range(1, n + 1):
        # number
        if str(i) == "1" * len(str(i)):
            dp[i] = str(i)
            dp_term[i] = str(i)
        for j in range(1, i):
            dp[i] = update(dp[i], dp[j] + "+" + dp[i - j])
            if i % j == 0 and j != 1:
                dp_term[i] = update(dp_term[i], dp_term[j] + "*" + dp_term[i // j])
        dp_term[i] = update(dp_term[i], "(" + dp[i] + ")")
        dp[i] = update(dp[i], dp_term[i])

    print(dp[n])
