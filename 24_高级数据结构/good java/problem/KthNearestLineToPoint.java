package template.problem;

import template.algo.DoubleBinarySearch;
import template.utils.GeoConstant;
import template.geometry.geo2.Circle2;
import template.geometry.geo2.Point2;
import template.primitve.generated.datastructure.DoubleArrayList;
import template.primitve.generated.datastructure.IntegerBIT;

import java.util.ArrayList;
import java.util.List;

public class KthNearestLineToPoint {
    /**
     * <pre>
     * 给定n个顶点，两两形成n(n-1)/2条直线，求这些直线中离center最近的第k条线离center的距离
     * </pre>
     * O(n\log_2n\log_2M)
     */
    public static double query(Point2 center, Point2[] points, long k, double relativeError, double absoluteError) {
        int n = points.length;
        Item[] items = new Item[n];
        double farthest = 0;
        for (int i = 0; i < n; i++) {
            items[i] = new Item();
            items[i].pt = points[i];
            farthest = Math.max(farthest, Point2.dist2(items[i].pt, center));
        }
        farthest = Math.sqrt(farthest) + 1e3;

        IntegerBIT bit = new IntegerBIT(2 * n);
        DoubleArrayList dal = new DoubleArrayList(2 * n);
        List<Item> outter = new ArrayList<>(n);
        DoubleBinarySearch dbs = new DoubleBinarySearch(relativeError, absoluteError) {
            @Override
            public boolean check(double mid) {
                Circle2 c = new Circle2(center, mid);
                List<Point2> a2 = new ArrayList<>(2);
                List<Point2> b2 = new ArrayList<>(2);
                outter.clear();
                dal.clear();
                bit.clear();
                for (Item p : items) {
                    a2.clear();
                    b2.clear();
                    if (c.contain(p.pt)) {
                        continue;
                    }
                    outter.add(p);
                    Circle2.tangent(c.center, c.r, p.pt, 0, false, a2, b2);
                    p.t1 = GeoConstant.theta(a2.get(0).x, a2.get(0).y);
                    p.t2 = GeoConstant.theta(a2.get(1).x, a2.get(1).y);
                    dal.add(p.t1);
                    dal.add(p.t2);
                }
                dal.unique();
                for (Item pt : outter) {
                    pt.l = dal.binarySearch(pt.t1) + 1;
                    pt.r = dal.binarySearch(pt.t2) + 1;
                    if (pt.l > pt.r) {
                        int tmp = pt.l;
                        pt.l = pt.r;
                        pt.r = tmp;
                    }
                }
                outter.sort((a, b) -> Integer.compare(a.r, b.r));
                long noIntersect = 0;
                for (Item pt : outter) {
                    noIntersect += bit.query(pt.l, pt.r);
                    bit.update(pt.r, 1);
                    bit.update(pt.l - 1, -1);
                }
                return (long) n * (n - 1) / 2 - noIntersect >= k;
            }
        };

        return dbs.binarySearch(0, farthest);
    }

    public static class Item {
        public Point2 pt;
        double t1;
        double t2;
        int l;
        int r;
    }
}

