package template.datastructure;

import template.math.DigitUtils;

public class SparseQuadTree implements Cloneable {
    SparseQuadTree tl, tr, bl, br;
    long area;

    public static SparseQuadTree build(int l, int r, int t, int b) {
        SparseQuadTree tree = new SparseQuadTree();
        tree.bl = tree.br = tree.tl = tree.tr = tree;
        return tree;
    }

    public static SparseQuadTree update(int ll, int rr, int bb, int tt, int l, int r, int b, int t, SparseQuadTree segment) {
        if (checkOutOfRange(ll, rr, bb, tt, l, r, b, t)) {
            return segment;
        }
        segment = segment.clone();
        if (checkCoverage(ll, rr, bb, tt, l, r, b, t)) {
            return segment;
        }

        int lrm = DigitUtils.floorAverage(l, r);
        int tbm = DigitUtils.floorAverage(t, b);

        segment.pushDown(l, r, b, t);
        segment.bl = update(ll, rr, bb, tt, l, lrm, b, tbm, segment.bl);
        segment.br = update(ll, rr, bb, tt, lrm + 1, r, b, tbm, segment.br);
        segment.tl = update(ll, rr, bb, tt, l, lrm, tbm + 1, t, segment.tl);
        segment.tr = update(ll, rr, bb, tt, lrm + 1, r, tbm + 1, t, segment.tr);
        segment.pushUp();

        return segment;
    }

    public static long query(int ll, int rr, int bb, int tt, int l, int r, int b, int t, SparseQuadTree segment) {
        if (checkOutOfRange(ll, rr, bb, tt, l, r, b, t)) {
            return 0;
        }
        if (checkCoverage(ll, rr, bb, tt, l, r, b, t)) {
            return segment.area;
        }

        int lrm = DigitUtils.floorAverage(l, r);
        int tbm = DigitUtils.floorAverage(t, b);

        segment.pushDown(l, r, b, t);
        return query(ll, rr, bb, tt, l, lrm, b, tbm, segment.bl) +
                query(ll, rr, bb, tt, lrm + 1, r, b, tbm, segment.br) +
                query(ll, rr, bb, tt, l, lrm, tbm + 1, t, segment.tl) +
                query(ll, rr, bb, tt, lrm + 1, r, tbm + 1, t, segment.tr);
    }

    public static boolean checkOutOfRange(int ll, int rr, int bb, int tt, int l, int r, int b, int t) {
        return ll > r || rr < l || tt < b || bb > t;
    }

    public static boolean checkCoverage(int ll, int rr, int bb, int tt, int l, int r, int b, int t) {
        return ll <= l && rr >= r && tt >= t && bb <= b;
    }

    public void pushUp() {
        area = tl.area + tr.area + bl.area + br.area;
    }


    public void pushDown(int l, int r, int b, int t) {
    }


    @Override
    public SparseQuadTree clone() {
        try {
            return (SparseQuadTree) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new RuntimeException(e);
        }
    }
}