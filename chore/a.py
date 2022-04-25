from collections import defaultdict
from typing import List


class ListNode:
    def __init__(self, x):
        self.val = x
        self.next = None


# 给一组链表，将其合并为环形链表（保证可以合并），链表中没有重复元素。然后从中间割裂开形成一条链，可以正着也可以反过来，要让这条链字典序最小。


#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# @param a ListNode类一维数组 指向每段碎片的开头
# @return ListNode类
#
class Solution:
    def solve(self, a: List[ListNode]):
        nexts = defaultdict(int)
        pre = defaultdict(int)
        visited = set()
        headP = a[0]
        while True:
            if headP.val in visited:
                break
            visited.add(headP.val)
            nexts[headP.val] = headP.next.val
            pre[headP.next.val] = headP.val
            headP = headP.next

        min_ = min(nexts.keys())
        pre_, next_ = pre[min_], nexts[min_]
        order = pre if pre_ < next_ else nexts
        n = len(nexts)
        dummy = ListNode(0)
        dummyP = dummy
        while n:
            dummyP.next = ListNode(min_)
            min_ = order[min_]
            dummyP = dummyP.next
            n -= 1
        return dummy.next


node1 = ListNode(1)
node2 = ListNode(2)
node3 = ListNode(3)
node4 = ListNode(4)
node1.next = node2
node2.next = node3
node3.next = node4
node4.next = node1
print(Solution().solve([node1, node2, node4]).__dict__)


node3 = ListNode(3)
node7 = ListNode(7)
node4 = ListNode(4)
node5 = ListNode(5)
node1 = ListNode(1)
node10 = ListNode(10)
node3.next = node7
node7.next = node4
node4.next = node5
node5.next = node1
node1.next = node10
node10.next = node3
res = Solution().solve([node3])
for _ in range(10):
    print(res.val)
    res = res.next
