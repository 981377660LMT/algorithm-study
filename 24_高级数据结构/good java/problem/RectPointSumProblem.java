package template.problem;

import template.math.DigitUtils;
import template.primitve.generated.datastructure.LongArrayList;
import template.primitve.generated.datastructure.LongBIT;
import template.utils.SequenceUtils;

import java.util.Arrays;
import java.util.Comparator;

public class RectPointSumProblem {
    private static int sign(int x) {
        return (Integer.bitCount(x) & 1) == 0 ? 1 : -1;
    }

    private static long interval(LongBIT bit, int l, int r) {
        if (l > r) {
            return 0;
        }
        return bit.query(r) - bit.query(l - 1);
    }

    public static long[] solve(Point2D[] pts, Query2D[] qs) {
        LongArrayList list = new LongArrayList(pts.length);
        for (Point2D pt : pts) {
            list.add(pt.y);
        }
        list.unique();
        Query2D[] sortByXl = qs.clone();
        Query2D[] sortByXr = qs.clone();
        Arrays.sort(sortByXl, Comparator.comparingLong(x -> x.xl));
        Arrays.sort(sortByXr, Comparator.comparingLong(x -> x.xr));
        Arrays.sort(pts, Comparator.comparingLong(x -> x.x));
        int xlIter = 0;
        int xrIter = 0;
        int ptIter = 0;
        int m = list.size();
        LongBIT bit = new LongBIT(m);
        while (xlIter < qs.length || xrIter < qs.length) {
            long min = Long.MAX_VALUE;
            if (xrIter < qs.length) {
                min = Math.min(sortByXr[xrIter].xr, min);
            }
            if (xlIter < qs.length) {
                min = Math.min(sortByXl[xlIter].xl - 1, min);
            }
            while (ptIter < pts.length && pts[ptIter].x <= min) {
                Point2D head = pts[ptIter++];
                bit.update(list.binarySearch(head.y) + 1, head.w);
            }
            while (xlIter < sortByXl.length && sortByXl[xlIter].xl - 1 == min) {
                Query2D head = sortByXl[xlIter++];
                head.ans -= interval(bit, list.lowerBound(head.yl) + 1,
                        list.upperBound(head.yr));
            }
            while (xrIter < sortByXr.length && sortByXr[xrIter].xr == min) {
                Query2D head = sortByXr[xrIter++];
                head.ans += interval(bit, list.lowerBound(head.yl) + 1,
                        list.upperBound(head.yr));
            }
        }
        return Arrays.stream(qs).mapToLong(x -> x.ans).toArray();
    }

    public static void prefixSum(Point2D[] pts) {
        Arrays.sort(pts, (a, b) -> {
            int cmp = Long.compare(a.x, b.x);
            if (cmp == 0) {
                cmp = Long.compare(a.y, b.y);
            }
            if (cmp == 0) {
                cmp = -Long.compare(a.w, b.w);
            }
            return cmp;
        });
        dac(pts, 0, pts.length - 1, new Point2D[pts.length]);
    }

    public static long[] solve(Point3D[] pts, Query3D[] qs) {
        int m = qs.length;
        Point3D[][] sub = new Point3D[4][m];
        long[] ans = new long[m];
        long[] xs = new long[2];
        long[] ys = new long[2];
        for (int i = 0; i < m; i++) {
            xs[0] = qs[i].xr;
            xs[1] = qs[i].xl - 1;
            ys[0] = qs[i].yr;
            ys[1] = qs[i].yl - 1;
            for (int j = 0; j < 2; j++) {
                for (int k = 0; k < 2; k++) {
                    sub[(j << 1) | k][i] = new Point3D(xs[j], ys[k], qs[i].zr, 0);
                    sub[(j << 1) | k][i].zl = qs[i].zl - 1;
                }
            }
        }
        prefixSum(SequenceUtils.pack(i -> new Point3D[i], sub[0], sub[1], sub[2], sub[3], pts));
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < 4; j++) {
                ans[i] += sign(j) * sub[j][i].sum;
            }
        }
        return ans;
    }

    public static void prefixSum(Point3D[] pts) {
        Arrays.sort(pts, (a, b) -> {
            int cmp = Long.compare(a.x, b.x);
            if (cmp == 0) {
                cmp = Long.compare(a.y, b.y);
            }
            if (cmp == 0) {
                cmp = Long.compare(a.z, b.z);
            }
            if (cmp == 0) {
                cmp = -Long.compare(a.w, b.w);
            }
            return cmp;
        });
        LongArrayList list = new LongArrayList(pts.length * 2);
        for (Point3D pt : pts) {
            list.add(pt.z);
            list.add(pt.zl);
        }
        list.unique();
        for (Point3D pt : pts) {
            pt.z = list.binarySearch(pt.z) + 1;
            pt.zl = list.binarySearch(pt.zl) + 1;
        }
        dac(pts, 0, pts.length - 1, new Point3D[pts.length], new LongBIT(list.size()));
        for (Point3D pt : pts) {
            pt.z = list.get((int) (pt.z - 1));
            pt.zl = list.get((int) (pt.zl - 1));
        }
    }

    private static void dac(Point3D[] pts, int l, int r, Point3D[] buf, LongBIT bit) {
        if (l == r) {
            pts[l].sum = pts[l].w;
            return;
        }
        int m = DigitUtils.floorAverage(l, r);
        dac(pts, l, m, buf, bit);
        dac(pts, m + 1, r, buf, bit);
        int i = l;
        int j = m + 1;
        int wpos = l;
        while (i <= m || j <= r) {
            if (j > r || i <= m && pts[i].y <= pts[j].y) {
                bit.update((int) pts[i].z, pts[i].w);
                buf[wpos++] = pts[i];
                i++;
            } else {
                pts[j].sum += bit.query((int) pts[j].z) - bit.query((int) pts[j].zl);
                buf[wpos++] = pts[j];
                j++;
            }
        }
        for (int t = l; t <= m; t++) {
            bit.update((int) pts[t].z, -pts[t].w);
        }
        assert wpos == r + 1;
        System.arraycopy(buf, l, pts, l, r - l + 1);
    }

    private static void dac(Point2D[] pts, int l, int r, Point2D[] buf) {
        if (l == r) {
            pts[l].sum = pts[l].w;
            return;
        }
        int m = DigitUtils.floorAverage(l, r);
        dac(pts, l, m, buf);
        dac(pts, m + 1, r, buf);
        int i = l;
        int j = m + 1;
        long sum = 0;
        int wpos = l;
        while (i <= m || j <= r) {
            if (j > r || i <= m && pts[i].y <= pts[j].y) {
                sum += pts[i].w;
                buf[wpos++] = pts[i];
                i++;
            } else {
                pts[j].sum += sum;
                buf[wpos++] = pts[j];
                j++;
            }
        }
        assert wpos == r + 1;
        System.arraycopy(buf, l, pts, l, r - l + 1);
    }

    public static class Point2D {
        public long x;
        public long y;

        public Point2D(long x, long y, long w) {
            this.x = x;
            this.y = y;
            this.w = w;
        }

        long sum;
        long w;

        @Override
        public String toString() {
            return "Point2D{" +
                    "x=" + x +
                    ", y=" + y +
                    ", w=" + w +
                    '}';
        }
    }

    public static class Point3D extends Point2D {
        public long z;
        private long zl = Long.MIN_VALUE;

        public Point3D(long x, long y, long z, long w) {
            super(x, y, w);
            this.z = z;
        }

        @Override
        public String toString() {
            return "Point3D{" +
                    "x=" + x +
                    ", y=" + y +
                    ", z=" + z +
                    ", w=" + w +
                    '}';
        }
    }

    public static class Query2D {
        public long xl;
        public long xr;
        public long yl;
        public long yr;
        private long ans;

        public Query2D(long xl, long xr, long yl, long yr) {
            assert xl <= xr && yl <= yr;
            this.xl = xl;
            this.xr = xr;
            this.yl = yl;
            this.yr = yr;
        }
    }


    public static class Query3D {
        public long xl;
        public long xr;
        public long yl;
        public long yr;
        public long zl;
        public long zr;

        public Query3D(long xl, long xr, long yl, long yr, long zl, long zr) {
            this.xl = xl;
            this.xr = xr;
            this.yl = yl;
            this.yr = yr;
            this.zl = zl;
            this.zr = zr;
        }

        @Override
        public String toString() {
            return "Query3D{" +
                    "xl=" + xl +
                    ", xr=" + xr +
                    ", yl=" + yl +
                    ", yr=" + yr +
                    ", zl=" + zl +
                    ", zr=" + zr +
                    '}';
        }
    }
}
