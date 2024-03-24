package template.string;

import java.util.ArrayDeque;
import java.util.Deque;

public class SuffixBalancedTree {
    private static final double FACTOR = 0.75;
    private static Node[] stk = new Node[0];
    private static int tail;
    public Node root;
    // 指针, ref
    private ObjectHolder<Node> objectHolder = new ObjectHolder<>();
    private Deque<Node> dq;

    private static class ObjectHolder<V> {
        V data;

        public void clear() {
            data = null;
        }
    }

    public SuffixBalancedTree(int cap) {
        dq = new ArrayDeque<>(cap + 1);
        root = Node.NIL;
        Node dummy = new Node(Integer.MIN_VALUE);
        dummy.next = dummy;
        dummy.occur = 0;
        dummy.offsetToTail = -1;
        dummy.weight = 0;
        dq.addFirst(dummy);
    }

    private boolean check() {
        collect(root);
        for (int i = 1; i < tail; i++) {
            if (stk[i - 1].weight >= stk[i].weight) {
                return false;
            }
            if (compare(stk[i - 1], stk[i]) >= 0) {
                return false;
            }
        }
        for (int i = 0; i < tail; i++) {
            if (stk[i].occur < 0 || stk[i].occur > 1) {
                return false;
            }
        }

        if (root.aliveSize + 1 != dq.size()) {
            return false;
        }

        return true;
    }


    public Node addPrefix(int x) {
        objectHolder.clear();
        root = insert(root, x, dq.peekFirst(), objectHolder, 0, 1);
        Node node = objectHolder.data;
        dq.addFirst(node);
        // assert check();
        return node;
    }

    public void removePrefix() {
        assert dq.size() > 1;
        Node deleted = dq.removeFirst();
        delete(root, deleted);
        // assert check();

        //clean or not
        if (root.aliveSize * 2 < root.size) {
            collect(root);
            int wpos = 0;
            for (int i = 0; i < tail; i++) {
                if (stk[i].occur == 0) {
                    continue;
                }
                stk[wpos++] = stk[i];
            }
            root = refactor(0, wpos - 1, 0, 1);
        }
    }

    public int rank(Node node) {
        return rank(root, node);
    }

    public int leq(IntSequence seq) {
        return rank(root, seq);
    }

    public Node sa(int k) {
        k++;
        return kth(root, k);
    }

    public int[] sa() {
        collect(root);
        int[] sa = new int[size()];
        int wpos = 0;
        for (int i = 0; i < tail; i++) {
            if (stk[i].occur == 0) {
                continue;
            }
            sa[wpos++] = stk[i].offsetToTail;
        }
        return sa;
    }

    public int size() {
        return root.aliveSize;
    }

    private static void ensureSpace(int n) {
        if (stk.length >= n) {
            return;
        }
        int nextSize = Math.max(1 << 16, stk.length);
        while (nextSize < n) {
            nextSize += nextSize;
        }
        stk = new Node[nextSize];
    }

    private static int insertCompare(Node a, int key, Node next) {
        if (a.key != key) {
            return Integer.compare(a.key, key);
        }
        return Double.compare(a.next.weight, next.weight);
    }

    private static int compare(Node root, IntSequence seq) {
        int len = seq.length();
        for (int i = 0; i < len; i++, root = root.next) {
            if (seq.get(i) != root.key) {
                return Integer.compare(root.key, seq.get(i));
            }
        }
        return 0;
    }

    private static int compare(Node a, Node b) {
        for (int i = 0; a != b; i++, a = a.next, b = b.next) {
            if (a.key != b.key) {
                return Integer.compare(a.key, b.key);
            }
        }
        return 0;
    }

    private static Node kth(Node root, int k) {
        if (root == Node.NIL) {
            return root;
        }
        root.pushDown();
        Node ans;
        if (root.left.aliveSize >= k) {
            ans = kth(root.left, k);
        } else if (root.left.aliveSize + root.occur >= k) {
            ans = root;
        } else {
            ans = kth(root.right, k - root.left.aliveSize - root.occur);
        }
        //push up for calc purpose
        root.pushUp();
        return ans;
    }

    private static int rank(Node root, IntSequence seq) {
        if (root == Node.NIL) {
            return 0;
        }

        int ans = 0;
//        root = refactor(root, L, R);
        root.pushDown();
        int compRes = compare(root, seq);
        if (compRes > 0) {
            ans += rank(root.left, seq);
        } else {
            ans += root.aliveSize - root.right.aliveSize;
            ans += rank(root.right, seq);
        }
//        root.pushUp();
        return ans;
    }

    private static int rank(Node root, Node node) {
        if (root == Node.NIL) {
            return 0;
        }
//        root = refactor(root, L, R);
        root.pushDown();
        int ans = 0;
        if (root == node) {
            ans += root.aliveSize - root.right.aliveSize;
        } else {
            int compRes = root.compareTo(node);
            if (compRes > 0) {
                ans += rank(root.left, node);
            } else {
                ans += root.aliveSize - root.right.aliveSize;
                ans += rank(root.right, node);
            }
        }
//        root.pushUp();
        return ans;
    }


    private static void init(int key, Node root, Node next, double weight) {
        root.key = key;
        root.weight = weight;
        root.next = next;
        root.occur++;
        root.offsetToTail = next.offsetToTail + 1;
        root.pushUp();
    }

    private static Node newNode(int key, Node next, double weight) {
        Node root = new Node();
        init(key, root, next, weight);
        return root;
    }

    private static Node insert(Node root, int key, Node next, ObjectHolder<Node> insertNode, double L, double R) {
        if (root == Node.NIL) {
            root = newNode(key, next, (L + R) / 2);
            insertNode.data = root;
            return root;
        }
        root.pushDown();
        int cmpRes = insertCompare(root, key, next);
        if (cmpRes == 0) {
            insertNode.data = root;
            init(key, root, next, root.weight);
        } else if (cmpRes > 0) {
            root.left = insert(root.left, key, next, insertNode, L, root.weight);
        } else {
            root.right = insert(root.right, key, next, insertNode, root.weight, R);
        }
        root.pushUp();
        root = refactor(root, L, R);
        return root;
    }

    private static void delete(Node root, Node node) {
        assert root != Node.NIL;
        root.pushDown();
        if (root == node) {
            root.occur--;
        } else {
            int compRes = root.compareTo(node);
            if (compRes > 0) {
                delete(root.left, node);
            } else {
                delete(root.right, node);
            }
        }
        root.pushUp();
    }

    private static void collect(Node root) {
        ensureSpace(root.size);
        tail = 0;
        collect0(root);
        assert tail == root.size;
    }


    private static Node refactor(Node root, double L, double R) {
        double threshold = root.size * FACTOR;
        if (root.left.size > threshold || root.right.size > threshold) {
            collect(root);
            root = refactor(0, tail - 1, L, R);
        }
        return root;
    }

    private static void collect0(Node root) {
        if (root == Node.NIL) {
            return;
        }
        root.pushDown();
        collect0(root.left);
        stk[tail++] = root;
        collect0(root.right);
    }

    private static Node refactor(int l, int r, double L, double R) {
        if (l > r) {
            return Node.NIL;
        }
        int m = (l + r) / 2;
        Node root = stk[m];
        root.weight = (L + R) / 2;
        root.left = refactor(l, m - 1, L, root.weight);
        root.right = refactor(m + 1, r, root.weight, R);
        root.pushUp();
        return root;
    }


    @Override
    public String toString() {
        collect(root);
        StringBuilder ans = new StringBuilder("{");
        for (int i = 0; i < tail; i++) {
            ans.append(stk[i]).append(',');
        }
        if (ans.length() > 1) {
            ans.setLength(ans.length() - 1);
        }
        ans.append("}");
        return ans.toString();
    }

    public static class Node implements Cloneable, Comparable<Node> {
        public static final Node NIL = new Node();

        Node left = NIL;
        Node right = NIL;
        int size;
        int aliveSize;
        int key;
        byte occur;
        public double weight;
        Node next;
        //prev means the floor node
        public int offsetToTail;

        static {
            NIL.left = NIL.right = NIL;
            NIL.size = NIL.aliveSize = 0;
            NIL.key = -1;
            NIL.offsetToTail = -1;
        }

        @Override
        public int compareTo(Node o) {
            return Double.compare(weight, o.weight);
        }

        public void pushUp() {
            if (this == NIL) {
                return;
            }
            size = left.size + right.size + 1;
            aliveSize = left.aliveSize + right.aliveSize + occur;
        }

        public void pushDown() {

        }

        private Node() {
        }

        private Node(int key) {
            this.key = key;
            pushUp();
        }

        @Override
        public String toString() {
            StringBuilder ans = new StringBuilder("[");
            int remain = 10;
            Node node = this;
            for (; node != null && remain > 0; node = node.next == node ? null : node.next, remain--) {
                ans.append(node.key).append(',');
            }
            if (node != null) {
                ans.append(",...,");
            }
            if (ans.length() > 1) {
                ans.setLength(ans.length() - 1);
            }
            ans.append("]");
            return ans.toString();
        }
    }
}
