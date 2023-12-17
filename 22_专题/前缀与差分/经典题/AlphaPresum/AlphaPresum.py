from typing import Callable, List, Union


def alphaPresum(
    sOrOrds: Union[str, List[int]], sigma=26, offset=97
) -> Callable[[int, int, int], int]:
    """
    给定字符集信息和字符s,返回一个查询函数.用于查询s[start:end]间char的个数.
    """
    if isinstance(sOrOrds, str):
        sOrOrds = [ord(v) for v in sOrOrds]
    preSum = [[0] * sigma for _ in range(len(sOrOrds) + 1)]
    for i in range(1, len(sOrOrds) + 1):
        preSum[i][:] = preSum[i - 1][:]
        preSum[i][sOrOrds[i - 1] - offset] += 1

    def getCountOfSlice(start: int, end: int, ord: int) -> int:
        if start < 0:
            start = 0
        if end > len(sOrOrds):
            end = len(sOrOrds)
        if start >= end:
            return 0
        return preSum[end][ord - offset] - preSum[start][ord - offset]

    return getCountOfSlice


if __name__ == "__main__":
    preSum = alphaPresum("abcdabcd")
    assert preSum(0, 2, 97) == 1
    assert preSum(0, 8, 97) == 2
