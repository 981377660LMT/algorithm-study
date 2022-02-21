from typing import List, Optional


class ListNode:
    # 注意题目给出的和评测的py文件里ListNode定义不一样，提交的时候要注释题目给出的ListNode
    # 因为内部judge有时候会判断返回的ListNode类是不是他自己的ListNode类
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


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
