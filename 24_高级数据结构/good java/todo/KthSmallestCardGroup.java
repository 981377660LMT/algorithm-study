package template.problem;

import template.rand.Randomized;

import java.util.Arrays;
import java.util.PriorityQueue;
import java.util.function.LongConsumer;

/**
 * There are n groups of cards, now select one card from each group, call such set valid.
 * You're supposed to find the k-th smallest sum set
 */
public class KthSmallestCardGroup {
    private long[][] groups;
    private int start;
    private long sum;

    public KthSmallestCardGroup(long[]... groups) {
        this.groups = groups;
        Arrays.sort(groups, (a, b) -> Integer.compare(a.length, b.length));
        start = 0;
        while (start < groups.length && groups[start].length == 1) {
            start++;
        }

        for (long[] g : groups) {
            Randomized.shuffle(g);
            Arrays.sort(g);
            for (int i = g.length - 1; i >= 1; i--) {
                g[i] -= g[i - 1];
            }
            sum += g[0];
        }
        Arrays.sort(groups, start, groups.length, (a, b) -> Long.compare(a[1], b[1]));
    }

    /**
     * <pre>
     * Find 1,2,...,kth smallest sum of valid set in order.
     * </pre>
     * <pre>
     * O(klog k) complexity
     * </pre>
     */
    public void theFirstKSmallestSet(int k, LongConsumer consumer) {
        if (start == groups.length) {
            consumer.accept(sum);
            return;
        }

        PriorityQueue<State> pq = new PriorityQueue<>(3 * k, (a, b) -> Long.compare(a.val, b.val));
        pq.add(new State(sum, start, 0));

        while (k > 0 && !pq.isEmpty()) {
            k--;
            State head = pq.remove();
            consumer.accept(head.val);

            if (head.cId + 1 < groups[head.gId].length) {
                pq.add(new State(head.val + groups[head.gId][head.cId + 1],
                        head.gId, head.cId + 1));
            }
            if (head.cId > 0 && head.gId + 1 < groups.length) {
                pq.add(new State(head.val + groups[head.gId + 1][1],
                        head.gId + 1, 1));

                if (head.cId == 1) {
                    pq.add(new State(head.val - groups[head.gId][1]
                            + groups[head.gId + 1][1], head.gId + 1, 1));
                }
            }
        }

    }

    private static class State {
        long val;
        int gId;
        int cId;

        public State(long val, int gId, int cId) {
            this.val = val;
            this.gId = gId;
            this.cId = cId;
        }

        @Override
        public String toString() {
            return String.format("(%d, %d, %d)", val, gId, cId);
        }
    }
}
