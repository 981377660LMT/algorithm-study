package template.algo;

import java.util.Arrays;

public class UndoXorDSU {
    int[] rank;
    int[] p;
    long[] xor;
    boolean conflict;

    public UndoXorDSU(int n) {
        rank = new int[n];
        p = new int[n];
        xor = new long[n];
    }

    public boolean conflict() {
        return conflict;
    }

    public void init() {
        Arrays.fill(rank, 1);
        Arrays.fill(p, -1);
        Arrays.fill(xor, 0);
        conflict = false;
    }

    public int find(int x) {
        while (p[x] != -1) {
            x = p[x];
        }
        return x;
    }

    public long xorToRoot(int x) {
        long ans = 0;
        while (p[x] != -1) {
            ans ^= xor[x];
            x = p[x];
        }
        return ans;
    }

    public long xor(int a, int b) {
        return xorToRoot(a) ^ xorToRoot(b);
    }

    public int size(int x) {
        return rank[find(x)];
    }

    public UndoOperation merge(int a, int b, long d) {
        return new UndoOperation() {
            int x, y;
            boolean conflictSnapshot;

            public void apply() {
                x = find(a);
                y = find(b);
                long delta = xorToRoot(a) ^ xorToRoot(b) ^ d;
                conflictSnapshot = conflict;
                if (x == y) {
                    conflict = conflict || delta != 0;
                    return;
                }
                if (rank[x] < rank[y]) {
                    int tmp = x;
                    x = y;
                    y = tmp;
                }
                p[y] = x;
                xor[y] = delta;
                rank[x] += rank[y];
            }


            public void undo() {
                int cur = y;
                while (p[cur] != -1) {
                    cur = p[cur];
                    rank[cur] -= rank[y];
                }
                p[y] = -1;
                xor[y] = 0;
                conflict = conflictSnapshot;
            }
        };
    }
}
