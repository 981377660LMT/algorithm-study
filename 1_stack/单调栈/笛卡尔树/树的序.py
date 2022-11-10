# 将一棵二叉查找树的键值插入序列称为树的生成序列，现给出一个生成序列，
# 求与其生成同样二叉查找树的所有生成序列中字典序最小的那个
# !n<=1e5 字典序最小的二叉搜索树的插入序列

# !给一个生成序列，建出一棵笛卡尔树，求字典序最小的可以得到相同笛卡尔树的生成序列
# !按题意建好树之后输出先序遍历即可,注意需要将索引和值调换

from typing import List
from 笛卡尔树 import buildCartesianTree2


def solve(insertNums: List[int]) -> List[int]:
    """字典序最小的二叉搜索树的插入序列

    Args:
        n (int): 长度
        nums (List[int]): 1-n 的排列
    """

    def preOrder(insertIndex: int) -> None:
        res.append(insertIndex)
        if leftChild[insertIndex] != -1:
            preOrder(leftChild[insertIndex])
        if rightChild[insertIndex] != -1:
            preOrder(rightChild[insertIndex])

    mp = {num: i for i, num in enumerate(insertNums, 1)}
    newNums = [mp[i] for i in range(1, len(insertNums) + 1)]  # !交换索引和值
    rootIndex, leftChild, rightChild = buildCartesianTree2(newNums)
    res = []
    preOrder(rootIndex)
    return [x + 1 for x in res]


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    print(*solve(nums))
