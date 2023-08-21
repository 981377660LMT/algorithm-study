# 2816. 翻倍以链表形式表示的数字
# https://leetcode.cn/problems/double-a-number-represented-as-a-linked-list/description/
# 给你一个 非空 链表的头节点 head ，表示一个不含前导零的非负数整数。
# 将链表 翻倍 后，返回头节点 head 。
# !链表中节点的数目在范围 [1, 1e4] 内

# 注意python3.10为了防止ddos攻击，将int的最大位数限制到了4300 (可以理解为python的爆int)
# ValueError: Exceeds the limit (4300) for integer string conversion: value has 4590 digits;
# use sys.set_int_max_str_digits() to increase the limit


import sys
from typing import Optional

sys.set_int_max_str_digits(0)


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


class Solution:
    def doubleIt(self, head: Optional[ListNode]) -> Optional[ListNode]:
        arr = linkedListToArray(head)
        num = int("".join(map(str, arr)))
        num *= 2
        arr = list(map(int, str(num)))
        return arrayToLinkedList(arr)
