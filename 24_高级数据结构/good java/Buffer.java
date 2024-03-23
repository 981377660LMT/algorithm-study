package template.utils;

import java.util.ArrayDeque;
import java.util.Deque;
import java.util.function.Consumer;
import java.util.function.Supplier;

public class Buffer<T> {
    private Deque<T> deque;
    private Supplier<T> supplier;
    private Consumer<T> cleaner;
    private int allocTime;
    private int releaseTime;


    public Buffer(Supplier<T> supplier) {
        this(supplier, (x) -> {
        });
    }

    public Buffer(Supplier<T> supplier, Consumer<T> cleaner) {
        this(supplier, cleaner, 0);
    }

    public Buffer(Supplier<T> supplier, Consumer<T> cleaner, int exp) {
        this.supplier = supplier;
        this.cleaner = cleaner;
        deque = new ArrayDeque<>(exp);
    }

    public T alloc() {
        allocTime++;
        T res;
        if (deque.isEmpty()) {
            res = supplier.get();
            cleaner.accept(res);
        } else {
            res = deque.removeFirst();
        }
        return res;
    }

    public void release(T e) {
        if (e == null) {
            return;
        }
        releaseTime++;
        cleaner.accept(e);
        deque.addLast(e);
    }

    public void release(T a, T b) {
        release(a);
        release(b);
    }

    public void release(T a, T b, T c) {
        release(a, b);
        release(c);
    }

    public void release(T a, T b, T c, T d) {
        release(a, b, c);
        release(d);
    }

    public void check() {
        if (allocTime != releaseTime) {
            throw new IllegalStateException("Buffer alloc " + allocTime + " but release " + releaseTime);
        }
    }
}
