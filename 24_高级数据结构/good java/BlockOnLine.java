package template.datastructure;

import java.util.Arrays;

public class BlockOnLine {
    int[] L;
    int[] R;
    int n;
    CallBack c;

    public static interface CallBack {
        void addBlock(int l, int r);

        void removeBlock(int l, int r);
    }

    public BlockOnLine(int n, CallBack c) {
        L = new int[n];
        R = new int[n];
        this.n = n;
        this.c = c;
        init();
    }

    public void init() {
        Arrays.fill(L, n);
        Arrays.fill(R, -1);
    }

    public void add(int i) {
        assert i >= 0 && i < n;
        assert L[i] > R[i];
        int from = i;
        int to = i;
        if (i > 0 && L[i - 1] <= R[i - 1]) {
            //cool
            from = L[i - 1];
            c.removeBlock(from, i - 1);
        }
        if (i + 1 < n && L[i + 1] <= R[i + 1]) {
            to = R[i + 1];
            c.removeBlock(i + 1, to);
        }
        L[from] = L[to] = from;
        R[from] = R[to] = to;
        c.addBlock(from, to);
    }
}
