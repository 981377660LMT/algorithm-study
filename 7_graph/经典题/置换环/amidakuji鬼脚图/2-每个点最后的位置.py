from typing import Dict, List


def amidakuji(n: int, lines: List[int]) -> Dict[int, int]:
    """
    鬼脚图中每个点最后的位置

    Args:
        n: 鬼脚图的n个出发点1-n
        lines: 鬼脚图的横线,一共m条,从上往下表示.
               lines[i] 表示第i条横线连接 line[i] 和 line[i]+1. (1-index)

    Returns:
        每个点最后的位置 (1-index)
    """
    starts = list(range(1, n + 1))
    for col in lines:
        starts[col - 1], starts[col] = starts[col], starts[col - 1]
    return {num: i for i, num in enumerate(starts, 1)}


assert amidakuji(5, [1, 4, 3, 4, 2, 3, 1]) == {5: 1, 2: 2, 4: 3, 1: 4, 3: 5}
