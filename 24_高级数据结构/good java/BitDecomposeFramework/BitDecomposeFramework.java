package template.algo;

import java.util.Iterator;
import java.util.function.BiFunction;

public class BitDecomposeFramework<E> implements Iterable<E> {
    private Object[] data;
    private BiFunction<E, E, E> merger;
    private int end;

    public BitDecomposeFramework(BiFunction<E, E, E> merger) {
        this.merger = merger;
        data = new Object[32];
    }

    public void add(E e) {
        add(e, 0);
    }

    private void add(E e, int index) {
        if (data[index] == null) {
            data[index] = e;
            end = Math.max(end, index + 1);
            return;
        }
        add(merger.apply((E) data[index], e), index + 1);
        data[index] = null;
    }

    @Override
    public Iterator<E> iterator() {
        return new Iterator<E>() {
            int index = 0;

            @Override
            public boolean hasNext() {
                return index < end;
            }

            @Override
            public E next() {
                while (index < end && data[index] == null) {
                    index++;
                }
                return (E) data[index++];
            }
        };
    }

    public void addAll(BitDecomposeFramework<E> x) {
        for (int i = 0; i < x.end; i++) {
            if (x.data[i] == null) {
                continue;
            }
            add((E) x.data[i], i);
        }
    }
}
