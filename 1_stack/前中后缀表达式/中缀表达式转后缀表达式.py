#
#  Convert infix notation to postfix notation
#
#  Description:
#  1 * 2 + 3 * (4 + 5) <- infix notation
#  1 2 * 3 4 5 + * +   <- postfix notation
#  Postfix notation is easy to evaluate.
#
#  Algorithm:
#  Shunting-yard algorithm by Dijkstra.
#
#  Verified:
#  SPOJ 4: ONP - Transform the Expression
#


END = "$"


def rank(c: str) -> int:
    return f"(^/*-+){END}".find(c)


def infixToPostFix(s: str) -> str:
    """中缀表达式转后缀表达式."""
    s += END  # terminal symbol
    sb = []
    op = [END]
    for c in s:
        if c.isalnum():
            sb.append(c)
            continue
        while op[-1] != "(":
            if rank(op[-1]) >= rank(c):
                break
            sb.append(op.pop())
        if c == ")":
            op.pop()
        else:
            op.append(c)
    return "".join(sb)


if __name__ == "__main__":
    print(infixToPostFix("1*2+3*(4+5)"))  # 12*345+*+
