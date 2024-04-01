package template.problem;

import template.datastructure.Range2DequeAdapter;
import template.datastructure.SegTree;
import template.datastructure.summary.ActiveSum;
import template.datastructure.summary.ActiveUpdate;
import template.geometry.geo2.IntegerRect2;
import template.primitve.generated.datastructure.LongArrayList;

import java.util.Arrays;
import java.util.Comparator;

public class RectUnionArea {
    /**
     * <pre>
     * 计算矩阵交的面积
     * 时间复杂度：O(n\log_2n)
     * 空间复杂度：O(n)
     * </pre>
     */
    public static long unionArea(IntegerRect2[] rects) {
        LongArrayList list = new LongArrayList(rects.length * 2);
        for (IntegerRect2 r : rects) {
            list.add(r.lb[0]);
            list.add(r.rt[0]);
        }
        list.unique();
        int m = list.size();
        if (m <= 1) {
            return 0;
        }
        for (IntegerRect2 rect : rects) {
            rect.lb[0] = list.binarySearch(rect.lb[0]);
            rect.rt[0] = list.binarySearch(rect.rt[0]);
        }
        IntegerRect2[] sortByB = rects.clone();
        IntegerRect2[] sortByT = rects.clone();
        Arrays.sort(sortByB, Comparator.comparingLong(x -> x.lb[1]));
        Arrays.sort(sortByT, Comparator.comparingLong(x -> x.rt[1]));
        Range2DequeAdapter<IntegerRect2> dqByB = new Range2DequeAdapter<>(i -> sortByB[i], 0, sortByB.length - 1);
        Range2DequeAdapter<IntegerRect2> dqByT = new Range2DequeAdapter<>(i -> sortByT[i], 0, sortByT.length - 1);
        SegTree<ActiveSum, ActiveUpdate> st = new SegTree<>(0, m - 2, ActiveSum::new,
                ActiveUpdate::new, i -> ActiveSum.ofUnactive(list.get(i + 1) - list.get(i)));
        ActiveUpdate bufUpd = new ActiveUpdate();
        long last = 0;
        long ans = 0;
        long total = list.get(m - 1) - list.get(0);
        while (!dqByT.isEmpty()) {
            long now = dqByT.peekFirst().rt[1];
            if (!dqByB.isEmpty()) {
                now = Math.min(now, dqByB.peekFirst().lb[1]);
            }
            ans += (total - st.sum.sumOfUnactiveCell()) * (now - last);
            while (!dqByB.isEmpty() && dqByB.peekFirst().lb[1] == now) {
                IntegerRect2 head = dqByB.removeFirst();
                bufUpd.asUpdate(1);
                st.update((int) head.lb[0], (int) head.rt[0] - 1, 0, m - 2, bufUpd);
            }
            while (!dqByT.isEmpty() && dqByT.peekFirst().rt[1] == now) {
                IntegerRect2 head = dqByT.removeFirst();
                bufUpd.asUpdate(-1);
                st.update((int) head.lb[0], (int) head.rt[0] - 1, 0, m - 2, bufUpd);
            }
            last = now;
        }
        for (IntegerRect2 rect : rects) {
            rect.lb[0] = list.get((int) rect.lb[0]);
            rect.rt[0] = list.get((int) rect.rt[0]);
        }
        return ans;
    }
}

