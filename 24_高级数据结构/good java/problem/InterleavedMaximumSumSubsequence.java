package template.problem;

import template.datastructure.SegTree;
import template.primitve.generated.datastructure.LongPreSum;
import template.utils.*;

import java.util.Arrays;
import java.util.function.Supplier;

/**
 * <p>
 * given array a with length n, define f(a, k) is the maximum weight of the interleaved subsequence.
 * </p>
 * an interleaved subsequence is a indices set, I[1], ..., I[k], that I[1] < I[2] < ... < I[k]
 * <p>
 * The weight of the interleaved subsequence is \sum_{i=1}^k (-1)^{i+1} a[I[i]]
 * </p>
 */
public class InterleavedMaximumSumSubsequence {
    public static interface Callback {
        public void callback(int size, long sum, boolean[] set);
    }

    static long inf = (long) 1e18;

    private static int sign(int x) {
        return (x & 1) == 1 ? -1 : 1;
    }

    /**
     * Find f(a, 0), f(a, 1), ... , f(a, a.length) and report by callback
     * @param a
     * @param cb
     */
    public static void solve(long[] a, Callback cb) {
        int n = a.length;
        SumImpl[] buf = new SumImpl[2 * n];
        for(int i = 0; i < buf.length; i++){
            buf[i] = new SumImpl();
        }
        PreAllocSupplier<SumImpl> pas = new PreAllocSupplier<>();
        pas.init(buf);

        //full
        SegTree<SumImpl, UpdateImpl> st = new SegTree<>(0, n - 1,
                pas, UpdateImpl::new, i -> {
            SumImpl s = pas.get();
            s.init(i, -sign(i) * a[i]);
            return s;
        });
        boolean[] set = new boolean[n];
        Arrays.fill(set, true);
        process(cb, st, n, set);

        LongPreSum[] ps = new LongPreSum[2];
        for (int i = 0; i < 2; i++) {
            int finalI = i;
            ps[i] = new LongPreSum(j -> sign(finalI ^ j) * a[j], n);
        }
        int delPos = 0;
        long best = -inf;
        for (int i = 0; i < n; i++) {
            long sum = ps[0].prefix(i - 1) +
                    ps[1].post(i + 1);
            if (sum > best) {
                best = sum;
                delPos = i;
            }
        }

        pas.init(buf);
        st.init(0, n - 1,
                SumImpl::new, UpdateImpl::new, i -> {
            SumImpl s = new SumImpl();
            s.init(i, -sign(i) * a[i]);
            return s;
        });
        u.clear();
        u.size = -1;
        st.update(delPos, delPos, 0, n - 1, u);
        u.clear();
        u.inv = 1;
        st.update(delPos, n - 1, 0, n - 1, u);
        Arrays.fill(set, true);
        set[delPos] = false;

        process(cb, st, n, set);
    }

    private static void process(Callback cb, SegTree<SumImpl, UpdateImpl> st, int n, boolean[] set) {
        while (st.sum.size >= 2) {
            SumImpl s = st.sum;
            cb.callback(st.sum.size, -s.sum[0], set);
            int l = s.minBegin[0];
            int r = s.minEnd[0];
            u.clear();
            u.size = -1;
            st.update(l, l, 0, n - 1, u);
            st.update(r, r, 0, n - 1, u);
            u.clear();
            u.inv = 1;
            st.update(l, r, 0, n - 1, u);
            set[l] = set[r] = false;
        }
        cb.callback(st.sum.size, -st.sum.sum[0], set);
    }

    static UpdateImpl u = new UpdateImpl();

    static class UpdateImpl extends CloneSupportObject<UpdateImpl> implements Update<UpdateImpl> {
        int inv;
        int size;

        @Override
        public void update(UpdateImpl update) {
            inv ^= update.inv;
            size += update.size;
        }

        @Override
        public void clear() {
            inv = 0;
            size = 0;
        }

        @Override
        public boolean ofBoolean() {
            return inv != 0 || size != 0;
        }
    }

    static class SumImpl implements UpdatableSum<SumImpl, UpdateImpl> {
        long[] sum = new long[2];
        long[] min = new long[2];
        long[] pref = new long[2];
        long[] suf = new long[2];
        int[] minBegin = new int[2];
        int[] minEnd = new int[2];
        int[] sufBegin = new int[2];
        int[] prefEnd = new int[2];
        int size;
        static long inf = (long) 1e18;

        public void init(int index, long x) {
            size = 1;
            sum[0] = x;
            sum[1] = -x;
            min[0] = min[1] = -inf;
            pref[0] = suf[0] = x;
            pref[1] = suf[1] = -x;
            prefEnd[0] = sufBegin[0] = prefEnd[1] = sufBegin[1] = index;
        }

        @Override
        public void add(SumImpl right) {
            if (this.size == 0) {
                copy(right);
                return;
            }
            if (right.size == 0) {
                return;
            }
            //min
            for (int i = 0; i < 2; i++) {
                if (min[i] < right.min[i]) {
                    min[i] = right.min[i];
                    minBegin[i] = right.minBegin[i];
                    minEnd[i] = right.minEnd[i];
                }
                long cand = suf[i] + right.pref[i];
                if (min[i] < cand) {
                    min[i] = cand;
                    minBegin[i] = sufBegin[i];
                    minEnd[i] = right.prefEnd[i];
                }
                cand = sum[i] * 2 + right.pref[i];
                if (pref[i] < cand) {
                    pref[i] = cand;
                    prefEnd[i] = right.prefEnd[i];
                }
                cand = suf[i] + right.sum[i] * 2;
                suf[i] = cand;
                if (suf[i] < right.suf[i]) {
                    suf[i] = right.suf[i];
                    sufBegin[i] = right.sufBegin[i];
                }
                sum[i] += right.sum[i];
            }
            size += right.size;
        }

        @Override
        public void copy(SumImpl right) {
            for (int i = 0; i < 2; i++) {
                sum[i] = right.sum[i];
                min[i] = right.min[i];
                pref[i] = right.pref[i];
                suf[i] = right.suf[i];
                minBegin[i] = right.minBegin[i];
                minEnd[i] = right.minEnd[i];
                prefEnd[i] = right.prefEnd[i];
                sufBegin[i] = right.sufBegin[i];
            }
            size = right.size;
        }

        @Override
        public SumImpl clone() {
            SumImpl ans = new SumImpl();
            ans.copy(this);
            return ans;
        }

        @Override
        public void update(UpdateImpl update) {
            if (update.inv == 1) {
                SequenceUtils.swap(sum, 0, 1);
                SequenceUtils.swap(min, 0, 1);
                SequenceUtils.swap(pref, 0, 1);
                SequenceUtils.swap(suf, 0, 1);
                SequenceUtils.swap(minBegin, 0, 1);
                SequenceUtils.swap(minEnd, 0, 1);
                SequenceUtils.swap(prefEnd, 0, 1);
                SequenceUtils.swap(sufBegin, 0, 1);
            }
            size += update.size;
        }

        @Override
        public String toString() {
            return "(" + size + "," + pref[0] + ")";
        }
    }
}
