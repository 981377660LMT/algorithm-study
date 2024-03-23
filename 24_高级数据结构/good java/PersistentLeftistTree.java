package template.graph;

import template.utils.CloneSupportObject;

import java.util.*;

// TODO: 学习接口、命名、设计
public class PersistentLeftistTree<K> extends CloneSupportObject<PersistentLeftistTree<K>> {
    public static final PersistentLeftistTree NIL = new PersistentLeftistTree<>(null);

    static {
        NIL.left = NIL;
        NIL.right = NIL;
        NIL.dist = -1;
    }

    PersistentLeftistTree<K> left = NIL;
    PersistentLeftistTree<K> right = NIL;
    int dist;
    K key;

    public PersistentLeftistTree(K key) {
        this.key = key;
    }

    public static <K> PersistentLeftistTree<K> createFromCollection(Collection<PersistentLeftistTree<K>> trees, Comparator<K> cmp) {
        return createFromDeque(new ArrayDeque<>(trees), cmp);
    }

    public static <K> PersistentLeftistTree<K> createFromDeque(Deque<PersistentLeftistTree<K>> deque, Comparator<K> cmp) {
        while (deque.size() > 1) {
            deque.addLast(merge(deque.removeFirst(), deque.removeFirst(), cmp));
        }
        return deque.removeLast();
    }

    public static <K> PersistentLeftistTree<K> merge(PersistentLeftistTree<K> a, PersistentLeftistTree<K> b, Comparator<K> cmp) {
        if (a == NIL) {
            return b;
        } else if (b == NIL) {
            return a;
        }
        if (cmp.compare(a.key, b.key) > 0) {
            PersistentLeftistTree<K> tmp = a;
            a = b;
            b = tmp;
        }
        a = a.clone();
        a.right = merge(a.right, b, cmp);
        if (a.left.dist < a.right.dist) {
            PersistentLeftistTree<K> tmp = a.left;
            a.left = a.right;
            a.right = tmp;
        }
        a.dist = a.right.dist + 1;
        return a;
    }

    public boolean isEmpty() {
        return this == NIL;
    }

    public K peek() {
        return key;
    }

    public static <K> PersistentLeftistTree<K> pop(PersistentLeftistTree<K> root, Comparator<K> cmp) {
        PersistentLeftistTree<K> ans = merge(root.left, root.right, cmp);
        return ans;
    }

    private void toStringDfs(StringBuilder builder) {
        if (this == NIL) {
            return;
        }
        builder.append(key).append(' ');
        left.toStringDfs(builder);
        right.toStringDfs(builder);
    }

    @Override
    public String toString() {
        StringBuilder builder = new StringBuilder();
        toStringDfs(builder);
        return builder.toString();
    }


    private static class PersistentLeftistTreeIteratorAdapter<V> implements Iterator<V> {
        private PersistentLeftistTree<V> tree;
        private Comparator<V> comparator;

        public PersistentLeftistTreeIteratorAdapter(PersistentLeftistTree<V> tree, Comparator<V> comparator) {
            this.tree = tree;
            this.comparator = comparator;
        }

        @Override
        public boolean hasNext() {
            return !tree.isEmpty();
        }

        @Override
        public V next() {
            V ans = tree.peek();
            tree = PersistentLeftistTree.pop(tree, comparator);
            return ans;
        }
    }

    public static <V> Iterator<V> asIterator(PersistentLeftistTree<V> heap, Comparator<V> comparator) {
        return new PersistentLeftistTreeIteratorAdapter<>(heap, comparator);
    }
}
