# 多个数组，多路归并
# 这道题思路跟丑数的思路差不多


from typing import List


class Solution:
    def arraysIntersection(self, arr1: List[int], arr2: List[int], arr3: List[int]) -> List[int]:
        """
        用三个指针同时在三数组上扫描：
        - 若三者相等，加入答案，并 i,j,k 同时 +1
        - 否则，把指向最小值的那个指针 +1
        时间：O(n1+n2+n3)，空间：O(min(n1,n2,n3))
        """
        i = j = k = 0
        res = []
        n1, n2, n3 = len(arr1), len(arr2), len(arr3)

        while i < n1 and j < n2 and k < n3:
            v1, v2, v3 = arr1[i], arr2[j], arr3[k]
            # 三者相等，找到一个交集元素
            if v1 == v2 == v3:
                res.append(v1)
                i += 1
                j += 1
                k += 1
            else:
                # 找到最小的那个，向前推进
                m = min(v1, v2, v3)
                if v1 == m:
                    i += 1
                if v2 == m:
                    j += 1
                if v3 == m:
                    k += 1

        return res


print(Solution().arraysIntersection([1, 2, 3, 4, 5], [1, 2, 5, 7, 9], [1, 3, 4, 5, 8]))
