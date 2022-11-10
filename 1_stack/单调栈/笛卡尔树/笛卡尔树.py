#  笛卡尔树的单调栈建树方法
# https://oi-wiki.org/ds/cartesian-tree/

# 给定一个1～n(n<=1e7)的排列p，构建其笛卡尔树。
# 即构建一棵二叉树，满足:
# !1.每个节点的编号满足二叉搜索树的性质。
# !2.节点主的权值为pi，每个节点的权值满足小根堆的性质。


from typing import List, Optional, Tuple


class Node:
    __slots__ = "weight", "key", "left", "right"

    def __init__(
        self, weight: int, key: int, left: Optional["Node"] = None, right: Optional["Node"] = None
    ):
        self.weight = weight
        """权值,满足小根堆的性质"""
        self.key = key
        """节点值,满足二叉搜索树的性质"""
        self.left = left
        self.right = right


# !指针版，返回根节点
def buildCartesianTree1(insertNums: List[int]) -> Optional["Node"]:
    n = len(insertNums)
    if n == 0:
        return None

    stack: List["Node"] = []  # !单增的单调栈维护`右链`
    for i, v in enumerate(insertNums):
        node = Node(i, v)
        last = None
        while stack and stack[-1].key > v:
            last = stack.pop()
        if stack:
            stack[-1].right = node  # !栈顶节点的值大于node的值,所以node是栈顶节点的右儿子
        if last is not None:
            node.left = last  # !栈里所有的值小于node的值,所以node的左儿子是栈底节点last
        stack.append(node)

    return stack[0]


# !非指针版，返回每个节点的左右儿子的编号
# !如果没有儿子，编号为-1
def buildCartesianTree2(insertNums: List[int]) -> Tuple[int, List[int], List[int]]:
    """笛卡尔树建树

    Args:
        nums (List[int]): 插入序列

    Returns:
        Tuple[int, List[int], List[int]]:
        根节点在插入序列中的索引,
        每个节点的左儿子在插入序列中的索引,
        每个节点的右儿子在插入序列中的索引
    """
    n = len(insertNums)
    leftChild, rightChild = [-1] * n, [-1] * n
    stack = []

    for i, v in enumerate(insertNums):
        last = -1
        while stack and insertNums[stack[-1]] > v:
            last = stack.pop()
        if stack:
            rightChild[stack[-1]] = i
        if last != -1:
            leftChild[i] = last
        stack.append(i)

    return stack[0], leftChild, rightChild


if __name__ == "__main__":
    # !如果笛卡尔树的 (key,weight) 键值对确定，且 key 互不相同，weight 互不相同，
    # !那么这个笛卡尔树的结构是唯一的

    def preOrder1(insertIndex: int) -> None:
        """前序遍历输出笛卡尔树结点插入顺序"""
        print(f"插入顺序key(满足BST):{insertIndex},结点weight(满足堆):{perm[insertIndex]}")
        if leftChild[insertIndex] != -1:
            preOrder1(leftChild[insertIndex])
        if rightChild[insertIndex] != -1:
            preOrder1(rightChild[insertIndex])

    perm = [9, 3, 7, 1, 8, 12, 10, 20, 15, 18, 5]
    rootIndex, leftChild, rightChild = buildCartesianTree2(perm)
    # preOrder1(rootIndex)

    def preOrder2(insertIndex: int) -> None:
        """前序遍历输出插入序列形成的的BST的节点值"""
        print(f"插入值key(满足BST):{allNums[insertIndex]},结点weight(满足堆):{newPerm[insertIndex]-1}")
        if leftChild2[insertIndex] != -1:
            preOrder2(leftChild2[insertIndex])
        if rightChild2[insertIndex] != -1:
            preOrder2(rightChild2[insertIndex])

    perm = [9, 3, 7, 1, 8, 12, 10, 20, 15, 18, 5]
    # !离散化到1-n
    allNums = sorted(set(perm))
    mapping = {num: i for i, num in enumerate(allNums, 1)}
    discretedPerm = [mapping[v] for v in perm]

    # !交换key和value
    mp = {num: i for i, num in enumerate(discretedPerm, 1)}
    newPerm = [mp[i] for i in range(1, len(mp) + 1)]
    rootIndex2, leftChild2, rightChild2 = buildCartesianTree2(newPerm)
    preOrder2(rootIndex2)
