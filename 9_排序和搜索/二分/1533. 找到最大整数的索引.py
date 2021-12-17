class ArrayReader(object):
    # Compares the sum of arr[l..r] with the sum of arr[x..y]
    # return 1 if sum(arr[l..r]) > sum(arr[x..y])
    # return 0 if sum(arr[l..r]) == sum(arr[x..y])
    # return -1 if sum(arr[l..r]) < sum(arr[x..y])
    def compareSub(self, l: int, r: int, x: int, y: int) -> int:
        ...

    # Returns the length of the array
    def length(self) -> int:
        ...


# arr 除了一个最大的整数外，其他所有整数都相等
class Solution:
    def getIndex(self, reader: 'ArrayReader') -> int:
        n = reader.length()
        left = 0
        right = n - 1

        while left < right:
            mid = left + right >> 1
            # -------------- 奇数个，L=0 [0 1 2] mid=3 [3 4 5 6] R=6
            if (right - left) % 2 == 0:
                # 比较 [0 1 2] 和 [4 5 6]
                check = reader.compareSub(left, mid - 1, mid + 1, right)
                # -- 左边大，在左
                if check == 1:
                    right = mid - 1  # 左半的R界为mid=2
                # -- 右边大，在右
                elif check == -1:
                    left = mid + 1  # 右半的L界为4
                # -- 两边一样大
                else:
                    return mid

            # --------------偶数个，L=0 [0 1 2] mid=2 [3 4 5] R=5
            else:
                # -- 比较 [0 1 2] 和 [3 4 5]
                check = reader.compareSub(left, mid, mid + 1, right)
                # -- 左边大，在左
                if check == 1:
                    right = mid  # 左半的R界为2
                # -- 右边大，在右
                elif check == -1:
                    left = mid + 1  # 右半的L界为3
                # -- 两边一样大，不存在
                else:
                    return -1
        return left

