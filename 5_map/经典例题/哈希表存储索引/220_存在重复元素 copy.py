import bisect


class Solution(object):
    def containsNearbyAlmostDuplicate(self, nums, k, t):
        from sortedcontainers import SortedSet

        st = SortedSet()
        left, right = 0, 0

        while right < len(nums):
            if right - left > k:
                st.remove(nums[left])
                left += 1
            index = bisect.bisect_left(st, nums[right] - t)
            if st and index >= 0 and index < len(st) and abs(st[index] - nums[right]) <= t:
                return True
            st.add(nums[right])
            right += 1
        return False

