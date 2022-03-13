from typing import List, Optional, Tuple

MOD = int(1e9 + 7)


class ListNode:
    def __init__(self, x):
        self.val = x
        self.next = None


class Solution:
    def isPalindrome(self, head: Optional['ListNode']) -> bool:
        def check(left, right, isMoved) -> bool:
            while left < right:
                if nums[left] != nums[right]:
                    if isMoved:
                        return False
                    else:
                        return check(left + 1, right, True) or check(left, right - 1, True)
                else:
                    left += 1
                    right -= 1
            return True

        nums = []
        while head:
            nums.append(head.val)
            head = head.next
        return check(0, len(nums) - 1, False)

