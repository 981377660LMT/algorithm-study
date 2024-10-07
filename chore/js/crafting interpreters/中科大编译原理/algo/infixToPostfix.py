# https://github.com/spaghetti-source/algorithm/blob/master/string/infix_to_postfix.cc
#
# # Convert infix notation to postfix notation
#
# Description:
#
# 	1 * 2 + 3 * (4 + 5) <- infix notation
# 	1 2 * 3 4 5 + * +   <- postfix notation
# 	Postfix notation is easy to evaluate.
#
# Algorithm:
#
# 	Shunting-yard algorithm by Dijkstra.
#
# Verified:
#
# 	SPOJ 4: ONP - Transform the Expression


def rank(c: str) -> int:
    if c == "(":
        return 0
    if c == "^":
        return 1
    if c == "/":
        return 2
    if c == "*":
        return 3
    if c == "-":
        return 4
    if c == "+":
        return 5
    if c == ")":
        return 6
    return 7


def infix_to_postfix(s: str) -> str:
    """中缀表达式转后缀表达式."""
    s += "_"  # sentinel
    sb = []
    op = ["_"]

    for c in s:
        if c.isalnum():
            sb.append(c)
        else:
            while op[-1] != "(":
                if rank(op[-1]) >= rank(c):
                    break
                sb.append(op[-1])
                op.pop()
            if c == ")":
                op.pop()
            else:
                op.append(c)
    return "".join(sb)


if __name__ == "__main__":
    print(infix_to_postfix("1*2+3*(4+5)"))  # 12*345+*+
