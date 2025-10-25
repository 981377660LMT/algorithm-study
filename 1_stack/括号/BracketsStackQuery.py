# https://atcoder.jp/contests/abc428/tasks/abc428_c
#
# 当一个字符串 T 满足以下条件时，我们称它为一个好的括号序列：
#
# 通过重复 0 次或多次以下操作，可以将 T 变为空字符串：
# 选择 T 中一个（连续的）子字符串 ()，并将其移除。
# 例如，()、(())() 以及空字符串都是好的括号序列，但 )( 和 ))) 不是。
# 现在有一个字符串 S，初始时为空。请按顺序处理 Q 个查询。在处理完每个查询后，你需要判断当前的 S 是否为一个好的括号序列。
#
# 查询有以下两种类型：
# 1 c: 给定一个字符 c，它要么是 ( 要么是 )。将 c 追加到 S 的末尾。
# 2: 删除 S 末尾的字符。保证在执行此查询时 S 不为空。
#
# !最终和为 0、所有前缀和非负(这里可以用StackAggregation查询栈的最小值)


class BracketsStackQuery:
    __slots__ = "_presum", "_premin"

    def __init__(self) -> None:
        self._presum = [0]
        self._premin = [0]

    def append(self, isLeftBracket: bool) -> None:
        v = 1 if isLeftBracket else -1
        self._presum.append(self._presum[-1] + v)
        self._premin.append(min(self._premin[-1], self._presum[-1]))

    def pop(self) -> None:
        self._presum.pop()
        self._premin.pop()

    def query(self) -> bool:
        return self._presum[-1] == 0 and self._premin[-1] == 0


if __name__ == "__main__":
    Q = int(input())
    stack = BracketsStackQuery()
    for _ in range(Q):
        query = input().split()
        if query[0] == "1":
            stack.append(query[1] == "(")
        else:
            stack.pop()
        print("Yes" if stack.query() else "No")
