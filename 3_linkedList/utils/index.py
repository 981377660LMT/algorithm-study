# 注意题目给出的和评测的py文件里ListNode定义不一样，
# 提交的时候要注释题目给出的ListNode
# 因为内部judge有时候会判断返回的ListNode类是不是他自己的ListNode类

from typing import List, Optional


class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


def arrayToLinkedList(nums):
    dummy = ListNode(0)
    p = dummy
    for num in nums:
        p.next = ListNode(num)
        p = p.next
    return dummy.next


def linkedListToArray(head):
    nums = []
    p = head
    while p:
        nums.append(p.val)
        p = p.next
    return nums
