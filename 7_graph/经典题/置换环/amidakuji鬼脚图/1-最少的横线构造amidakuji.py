# 冒泡排序邻位交换 构造amidakuji的横线
from typing import List


def amidakuji(ends: List[int]) -> List[int]:
    """
    冒泡排序构造amidakuji的横线

    Args:
        ends: 每个点最后的位置 (1-index)
    Returns:
        鬼脚图的横线,一共m条,从上往下表示.
        lines[i] 表示第i条横线连接 line[i] 和 line[i]+1. (1-index)
    """
    n = len(ends)
    res = []
    for i in range(n - 1):
        isSorted = True
        for j in range(n - 1 - i):
            if ends[j] > ends[j + 1]:
                isSorted = False
                res.append(j + 1)
                ends[j], ends[j + 1] = ends[j + 1], ends[j]
        if isSorted:
            break
    return res[::-1]  # 目标数组冒泡排序的逆序才是从上往下的横线顺序


assert amidakuji([5, 2, 4, 1, 3]) == [1, 3, 2, 4, 3, 2, 1]
