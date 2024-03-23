package template.utils;

import template.binary.Log2;

public class GridIndex {
    int mask;
    int bit;

    public GridIndex(int n, int m) {
        assert (long) n * m <= Integer.MAX_VALUE && m > 0 && n > 0;
        bit = Log2.ceilLog(m);
        mask = (1 << bit) - 1;
    }

    public int encode(int x, int y) {
        return (x << bit) | y;
    }

    public int decodeX(int id) {
        return id >>> bit;
    }

    public int decodeY(int id) {
        return id & mask;
    }
}
