# Definition for singly-linked list.
from typing import List, Optional

# 题目给出的和评测的py文件里ListNode的定义不一样，提交的时候要注释题目给出的ListNode
# 题目给的ListNode定义和他内部judge时候用的ListNode定义不一样
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next


class Solution:
    def mergeNodes(self, head: Optional['ListNode']) -> Optional['ListNode']:
        nums = []
        p = head
        while p:
            nums.append(int(p.val))
            p = p.next

        newNums = []
        curSum = 0
        for num in nums:
            curSum += num
            if num == 0:
                newNums.append(curSum)
                curSum = 0
        newNums = [num for num in newNums if num != 0]

        res = ListNode(0)
        p = res
        for num in newNums:
            p.next = ListNode(num)
            p = p.next
        return res.next


def makeLinkedList(nums: List[int]) -> Optional[ListNode]:
    head = ListNode(0)
    p = head
    for num in nums:
        p.next = ListNode(num)
        p = p.next
    return head.next


def printLinkedList(head: Optional[ListNode]) -> None:
    p = head
    while p:
        print(p.val, end=" ")
        p = p.next
    print()


node = makeLinkedList([0, 3, 1, 0, 4, 5, 2, 0])
printLinkedList(Solution().mergeNodes(node))

node = makeLinkedList([0, 1, 0, 3, 0, 2, 2, 0])
printLinkedList(Solution().mergeNodes(node))
