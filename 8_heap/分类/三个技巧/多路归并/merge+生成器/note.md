**多路归并求第 k 个...对**的问题

- `多路归并，每一路都是一个生成器(递增的数据流)`
- `heap.merge(*allGen) 相当于自动执行多路归并`
- `itertools.islice(iterable, k)相当于取出多路的前k个结果`

例如
merge(
(1, 4, 7,...), # 第一路(生成器)
(2, 5, 8,...), # 第二路(生成器)
(3, 6, 9,...), # 第三路(生成器)
)

取切片(islice)就是每次从每路取一个元素

```Python
def mergeTwo(nums1: List[int], nums2: List[int], k: int) -> List[int]:
    """两个有序数组选前k小的和"""
    gen = lambda index: (nums1[index] + num for num in nums2)  # 递增的一路
    allGen = [gen(i) for i in range(len(nums1))]  # 多路
    iterable = merge(*allGen)  # merge 相当于多路归并
    return list(islice(iterable, k))


class Solution:
    def kthSmallest(self, mat: List[List[int]], k: int) -> int:
        return next(reversed(reduce(lambda row1, row2: mergeTwo(row1, row2, k), mat)))
```
