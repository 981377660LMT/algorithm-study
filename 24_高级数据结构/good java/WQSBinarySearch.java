package template.algo;

public abstract class WQSBinarySearch {
    protected abstract double getBest();

    protected abstract int getTime();

    protected abstract void check(double costPerOperation);

    /**
     * top decide whether the charts of f(x)  look like a up convex
     * while x is the operation time and f(x) is the best profit.
     */
    public double search(double l, double r, int round, int k, boolean top) {
        if (l > r || round <= 0) {
            throw new IllegalArgumentException();
        }
        double m = 0;
        while (round-- > 0) {
            m = (l + r) / 2;
            check(m);
            if (getTime() == k) {
                break;
            }
            if (getTime() > k == top) {
                l = m;
            } else {
                r = m;
            }
        }
        return getBest() + m * k;
    }
}
