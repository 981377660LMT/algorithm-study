package template.problem;

import template.datastructure.BitSet;
import template.math.GCDs;
import template.math.LCMs;
import template.math.LongMillerRabin;
import template.primitve.generated.datastructure.LongArrayList;

import java.util.stream.LongStream;

public class MaximizingLCM {
    /**
     * <p>
     * Choose at most N number from [1, M], let's denote them as x[1],...,x[N]
     * maximize lcm(x[1], ..., x[N]).
     * </p>
     * <p>
     * time complexity: very fast, lots of optimizations
     * </p>
     *
     * @param N
     * @param M
     * @return chosen set x
     */
    public long[] maximize(int N, long M) {
        N = Math.min(N, 62);
        if (N >= M) {
            return LongStream.range(1, M + 1).toArray();
        }
        LongArrayList candList = new LongArrayList();
        int prime = 0;
        for (long i = M; i >= 1 && prime < N; i--) {
            candList.add(i);
            if (LongMillerRabin.mr(i, 10)) {
                prime++;
            }
        }

        bestVal = 0;
        cand = candList.toArray();
        if (prime < N) {
            LongArrayList filteredCand = new LongArrayList();
            for (long c : cand) {
                boolean ok = true;
                for (long t : cand) {
                    if (t > c && t % c == 0) {
                        ok = false;
                    }
                }
                if (ok) {
                    filteredCand.add(c);
                }
            }
            cand = filteredCand.toArray();
            added = new BitSet(cand.length);
            bestSet = new BitSet(cand.length);
            bf(N, 0, 1);
        } else {

            int k = cand.length;
            allow = new BitSet[N + 1];
            added = new BitSet(k);
            bestSet = new BitSet(k);
            conflict = new BitSet[k];
            for (int i = 0; i <= N; i++) {
                allow[i] = new BitSet(k);
            }
            for (int i = 0; i < k; i++) {
                conflict[i] = new BitSet(k);
                conflict[i].fill(true);
                for (int j = 0; j < k; j++) {
                    if (GCDs.gcd(cand[i], cand[j]) != 1) {
                        conflict[i].clear(j);
                    }
                }
            }
            allow[N].fill(true);
            smart(N, 0, 1);
        }

        LongArrayList sol = new LongArrayList(bestSet.size());
        for (int i = 0; i < cand.length; i++) {
            if (bestSet.get(i)) {
                sol.add(cand[i]);
            }
        }
        return sol.toArray();
    }

    long[] cand;
    BitSet[] allow;
    BitSet[] conflict;
    BitSet added;
    long bestVal;
    BitSet bestSet;

    public void updateBest(long bestVal, BitSet set) {
        if (bestVal > this.bestVal) {
            this.bestVal = bestVal;
            bestSet.copy(set);
        }
    }


    public void bf(int remain, int startIndex, long val) {
        if (remain == 0 || startIndex == cand.length) {
            updateBest(val, added);
            return;
        }
        bf(remain, startIndex + 1, val);
        added.set(startIndex);
        bf(remain - 1, startIndex + 1, LCMs.lcm(val, cand[startIndex]));
        added.clear(startIndex);
    }


    public void smart(int remain, int startIndex, long val) {
        if (remain == 0 || allow[remain].isEmpty()) {
            updateBest(val, added);
            return;
        }
        for (int bit = allow[remain].nextSetBit(startIndex); bit < allow[remain].capacity(); bit = allow[remain].nextSetBit(bit + 1)) {
            allow[remain - 1].copy(allow[remain]);
            allow[remain - 1].and(conflict[bit]);
            added.set(bit);
            smart(remain - 1, bit + 1, val * cand[bit]);
            added.clear(bit);
        }
    }
}
