package template.datastructure;

import java.util.AbstractCollection;
import java.util.Comparator;
import java.util.Iterator;
import java.util.PriorityQueue;

public class FixedMinCollection<V> extends AbstractCollection<V> {
    private PriorityQueue<V> pq;
    private int cap;
    private Comparator<V> comp;

    public FixedMinCollection(int cap, Comparator<V> comp) {
        if (cap == 0) {
            throw new IllegalArgumentException();
        }
        this.cap = cap;
        this.comp = comp;
        pq = new PriorityQueue<>(cap, comp.reversed());
    }

    @Override
    public boolean add(V v) {
        if (pq.size() < cap) {
            pq.add(v);
            return true;
        }
        if (comp.compare(pq.peek(), v) > 0) {
            pq.remove();
            pq.add(v);
            return true;
        }
        return false;
    }

    public V peekMax() {
        return pq.peek();
    }

    public V popMax() {
        return pq.remove();
    }

    @Override
    public void clear() {
        pq.clear();
    }

    @Override
    public Iterator<V> iterator() {
        return pq.iterator();
    }

    @Override
    public int size() {
        return pq.size();
    }
}
