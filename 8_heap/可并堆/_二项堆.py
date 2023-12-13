# https://zhuanlan.zhihu.com/p/98859548
# https://blog.51cto.com/u_13393656/3065850
# ①make-heap()：创建空堆
# ②insert(H,x)：结点x插入到堆H中
# ③minimum(H)：取最小值
# ④extract-min(H)：抽取最小值
# ⑤union(H1,H2)：合并堆H1，H2
# ⑥decrease-key(H,x,k)：x减值为k
# ⑦delete(H,x)：从堆H中删除结点x

# https://blog.csdn.net/yzf0011/article/details/56004750


class Node:
    __slots__ = ("key", "degree", "parent", "child", "sibling")

    def __init__(self, key):
        self.key = key
        self.degree = 0
        self.parent = None
        self.child = None
        self.sibling = None


class BinomialHeap:
    def __init__(self):
        self.head = None

    def make_heap(self):
        self.head = None

    def insert(self, key):
        node = Node(key)
        temp = BinomialHeap()
        temp.head = node
        self.union(temp)

    def minimum(self):
        if self.head is None:
            return None
        y = None
        x = self.head
        min_val = x.key
        while x is not None:
            if x.key < min_val:
                min_val = x.key
                y = x
            x = x.sibling
        return y

    def extract_min(self):
        if self.head is None:
            return None
        min_node_prev = None
        min_node = self.head
        cur_node = self.head
        next_node = cur_node.sibling
        while next_node is not None:
            if next_node.key < min_node.key:
                min_node_prev = cur_node
                min_node = next_node
            cur_node = next_node
            next_node = cur_node.sibling
        if min_node_prev is None:
            self.head = min_node.sibling
        else:
            min_node_prev.sibling = min_node.sibling
        child = None
        temp = min_node.child
        while temp is not None:
            sibling = temp.sibling
            temp.sibling = child
            child = temp
            temp = sibling
        temp_heap = BinomialHeap()
        temp_heap.head = child
        self.union(temp_heap)
        return min_node.key

    def union(self, heap):
        self.merge(heap)
        if self.head is None:
            return
        prev_x = None
        x = self.head
        next_x = x.sibling
        while next_x is not None:
            if x.degree != next_x.degree or (
                next_x.sibling is not None and next_x.sibling.degree == x.degree
            ):
                prev_x = x
                x = next_x
            elif x.key <= next_x.key:
                x.sibling = next_x.sibling
                self._link(next_x, x)
            else:
                if prev_x is None:
                    self.head = next_x
                else:
                    prev_x.sibling = next_x
                self._link(x, next_x)
                x = next_x
            next_x = x.sibling

    def decrease_key(self, node, new_key):
        if node is None or node.key < new_key:
            return
        node.key = new_key
        y = node
        z = y.parent
        while z is not None and y.key < z.key:
            y.key, z.key = z.key, y.key
            y = z
            z = y.parent

    def delete(self, node):
        if node is None or self.head is None:
            return
        self.decrease_key(node, float("-inf"))
        self.extract_min()

    def merge(self, heap):
        self.head = self._merge(self.head, heap.head)

    def _link(self, y, z):
        y.parent = z
        y.sibling = z.child
        z.child = y
        z.degree += 1

    def _merge(self, h1, h2):
        if h1 is None:
            return h2
        if h2 is None:
            return h1
        if h1.degree < h2.degree:
            h1.sibling = self._merge(h1.sibling, h2)
            return h1
        else:
            h2.sibling = self._merge(h1, h2.sibling)
            return h2


if __name__ == "__main__":
    pq = BinomialHeap()
    pq.insert(3)
    pq.insert(1)
    pq.insert(2)
    print(pq.extract_min())
    print(pq.extract_min())
    print(pq.extract_min())
    print(pq.extract_min())
