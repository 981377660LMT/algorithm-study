'''
舞蹈链（双向十字链表）
应用场景：精确覆盖（数独，智慧珠）
https://zhuanlan.zhihu.com/p/67324277
'''

# https://github.com/boristown/leetcode/blob/main/DLX/DLX.py
from typing import List


class CrossCycleLinkNode(object):
    def __init__(self, x, row):
        self.val = x
        self.row = row
        self.col = self
        self.up = self
        self.down = self
        self.left = self
        self.right = self

    def removeColumn(self):
        node = self
        while True:
            node.left.right = node.right
            node.right.left = node.left
            node = node.down
            if node == self:
                break

    def restoreColumn(self):
        node = self
        while True:
            node.left.right = node
            node.right.left = node
            node = node.down
            if node == self:
                break

    def removeRow(self):
        node = self
        while True:
            node.up.down = node.down
            node.down.up = node.up
            node = node.right
            if node == self:
                break

    def restoreRow(self):
        node = self
        while True:
            node.up.down = node
            node.down.up = node
            node = node.right
            if node == self:
                break


# head是头结点，res是装答案（行号）的列表
def dance_link_x(head: CrossCycleLinkNode, res: List):
    if head.right == head:
        return True

    node = head.right
    while node != head:
        if node.down == node:
            return False
        node = node.right

    restores = []
    first_col = head.right
    first_col.removeColumn()
    restores.append(first_col.restoreColumn)

    node = first_col.down
    while node != first_col:
        if node.right != node:
            node.right.removeRow()
            restores.append(node.right.restoreRow)
        node = node.down

    cur_restores_count = len(restores)
    selected_row = first_col.down
    while selected_row != first_col:
        res.append(selected_row.row)
        if selected_row.right != selected_row:
            row_node = selected_row.right
            while True:
                col_node = row_node.col
                col_node.removeColumn()
                restores.append(col_node.restoreColumn)
                col_node = col_node.down
                while col_node != col_node.col:
                    if col_node.right != col_node:
                        col_node.right.removeRow()
                        restores.append(col_node.right.restoreRow)
                    col_node = col_node.down
                row_node = row_node.right
                if row_node == selected_row.right:
                    break

        if dance_link_x(head, res):
            # while len(restores): restores.pop()()
            return True

        res.pop()
        while len(restores) > cur_restores_count:
            restores.pop()()
        selected_row = selected_row.down

    while len(restores):
        restores.pop()()
    return False


def initCol(col_count: int):
    head = CrossCycleLinkNode('head', 'column')
    for i in range(col_count):
        col_node = CrossCycleLinkNode(x=i, row=head.row)
        col_node.right = head
        col_node.left = head.left
        col_node.right.left = col_node
        col_node.left.right = col_node
    return head


def appendRow(head: CrossCycleLinkNode, row_id: int, nums: List):
    last = None
    col = head.right
    for num in nums:
        while col != head:
            if col.val == num:
                node = CrossCycleLinkNode(1, row_id)
                node.col = col
                node.down = col
                node.up = col.up
                node.down.up = node
                node.up.down = node
                if last is not None:
                    node.left = last
                    node.right = last.right
                    node.left.right = node
                    node.right.left = node
                last = node
                break
            col = col.right
        else:
            return
