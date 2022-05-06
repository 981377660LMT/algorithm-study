# 将数组变为严格递增的最小操作数
# 每次操作可以改变任一个数字

from typing import List
from LIS模板 import LIS

# 先变为非严格递增，就变为了LIS长度
def minRemove(nums: List[int]) -> int:
    nums = [num - i for num, i in enumerate(nums)]
    return len(nums) - LIS(nums, isStrict=False)


if __name__ == '__main__':
    assert (minRemove([1, 2, 6, 5, 4])) == 2

