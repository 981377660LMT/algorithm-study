# from typing import List


from typing import List


def heapSort(nums: List[int]) -> None:
    """
    堆排序
    先heapify,然后不断`交换堆顶和最后一个元素,并pushDown(0,size)`
    """

    def pushDown(index: int, size: int) -> None:
        """
        将nums[start]和nums[start*2+1]和nums[start*2+2]三个元素中最小的元素交换到nums[start]
        """
        root, left, right = index, (index << 1) + 1, (index << 1) + 2
        smallest = root
        if left < size and nums[left] < nums[smallest]:
            smallest = left
        if right < size and nums[right] < nums[smallest]:
            smallest = right
        if smallest != root:
            nums[root], nums[smallest] = nums[smallest], nums[root]
            pushDown(smallest, size)

    def heapify() -> None:
        """heapify最小堆化
        """
        last = len(nums) - 1
        parent = (last - 1) >> 1
        for i in range(parent, -1, -1):
            pushDown(i, size=len(nums))

    heapify()
    for i in range(len(nums) - 1, 0, -1):
        nums[i], nums[0] = nums[0], nums[i]
        pushDown(0, size=i)


if __name__ == '__main__':
    nums = [2, 3, 4, 1, 5, 4, 7]
    heapSort(nums)
    print(nums)

