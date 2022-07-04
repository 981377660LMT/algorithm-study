# n<=2e5
# 从前序遍历和中序遍历还原二叉树
# 如果存在 输出每个结点左右儿子 (空为0)
# 如果不存在 输出-1
from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def main() -> None:
    def dfs(pLeft: int, pRight: int, iLeft: int, iRight: int) -> int:
        """
        不要每次都暴力index从头查找根节点出现在中序遍历哪个位置
        因为元素全都不同 要用哈希表记录查找
        """
        if pLeft > pRight or iLeft > iRight:
            return 0

        root = preOrder[pLeft]
        rootIndex = mapping[root]
        if not (iLeft <= rootIndex <= iRight):
            raise ValueError("rootIndex not in range")

        leftLen, rightLen = rootIndex - iLeft, iRight - rootIndex
        leftSon[root] = dfs(pLeft + 1, pLeft + leftLen, iLeft, rootIndex - 1)
        rightSon[root] = dfs(pRight - rightLen + 1, pRight, rootIndex + 1, iRight)
        return root

    n = int(input())
    preOrder = [0] + list(map(int, input().split()))
    inOrder = [0] + list(map(int, input().split()))

    # !注意根节点必须要是1
    if preOrder[1] != 1:
        print(-1)
        exit(0)

    leftSon, rightSon = defaultdict(int), defaultdict(int)
    mapping = {v: i for i, v in enumerate(inOrder[1:], start=1)}
    try:
        dfs(1, n, 1, n)
    except ValueError:
        print(-1)
        exit(0)

    for i in range(1, n + 1):
        print(leftSon[i], rightSon[i])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            try:
                main()
            except (EOFError, ValueError):
                break
    else:
        main()
