class PolyNode:
    def __init__(self, x=0, y=0, next=None):
        self.coefficient = x
        self.power = y
        self.next = next


class Solution:
    def addPoly(self, poly1: 'PolyNode', poly2: 'PolyNode') -> 'PolyNode':
        dummy = PolyNode()
        dummyP = dummy

        while poly1 or poly2:
            if not poly2 or (poly1 and poly1.power > poly2.power):
                dummyP.next = poly1
                poly1 = poly1.next
                dummyP = dummyP.next
            elif not poly1 or (poly1 and poly2.power > poly1.power):
                dummyP.next = poly2
                poly2 = poly2.next
                dummyP = dummyP.next
            else:
                c = poly1.coefficient + poly2.coefficient
                if c != 0:
                    dummyP.next = PolyNode(c, poly1.power)
                    dummyP = dummyP.next
                poly1 = poly1.next
                poly2 = poly2.next

        dummyP.next = None
        return dummy.next

