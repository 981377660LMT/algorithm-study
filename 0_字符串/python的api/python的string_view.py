# python的string_view
btArr = bytearray([1, 2, 3, 2, 3, 6, 7, 8, 9, 10])

print(memoryview(btArr))
s1, s2 = btArr[1:3], btArr[3:5]
print(s1 == s2)
print(bytes(btArr[1:3]))
# str和bytearray的切片操作会产生新的切片str和bytearry并拷贝数据，使用memoryview之后不会。

##########################################################################################
# 2261. 含最多 K 个可整除元素的子数组
from typing import List


class Solution:
    def countDistinct(self, nums: List[int], k: int, p: int) -> int:
        res = set()
        arr = memoryview(bytearray(nums))
        for start in range(len(nums)):
            count = 0
            for end in range(start, len(nums)):
                if nums[end] % p == 0:
                    count += 1

                if count <= k:
                    res.add(bytes(arr[start : end + 1]))
                else:
                    break

        return len(res)


##########################################################################################

s = '1234u'
arr = memoryview(bytearray(map(ord, s)))
print(arr)
