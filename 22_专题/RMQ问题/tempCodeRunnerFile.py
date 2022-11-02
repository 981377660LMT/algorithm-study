ass Solution:
    def closestToTarget(self, arr: List[int], target: int) -> int:
        """
        1. 静态区间查询使用st表
        st表适用于区间重复贡献的问题
        时间复杂度O(nlog(n))
        2. 与运算具有单调性，可以使用二分查找
        """

        st = SparseTable(arr, and_)
        res = abs(arr[0] - target)
        for start in range(len(arr)):
            left, right = start, len(arr) - 1
            while left <= right:
                mid = (left + right) // 2
                # 越往左越大 越往右越小
                diff = st.query(start, mid) - target
                res = min(res, abs(diff))
                if diff == 0:
                    return 0
                elif diff > 0:
                    left = mid + 1
                else:
                    right = mid - 1
        return res
