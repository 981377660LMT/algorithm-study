package template.algo;

import template.math.DigitUtils;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public class OfflineUndoSegment {
    public interface Callback {
        void accept(int index);
    }

    private OfflineUndoSegment left;
    private OfflineUndoSegment right;
    private List<FlaggedCommutativeUndoOperation> items = new ArrayList<>();

    private void modify(FlaggedCommutativeUndoOperation item) {
        items.add(item);
    }

    public void pushUp() {
    }

    public void pushDown() {
    }

    public OfflineUndoSegment(int l, int r) {
        if (l < r) {
            int m = DigitUtils.floorAverage(l, r);
            left = new OfflineUndoSegment(l, m);
            right = new OfflineUndoSegment(m + 1, r);
            pushUp();
        } else {

        }
    }

    private boolean enter(int ll, int rr, int l, int r) {
        return ll <= l && rr >= r;
    }

    private boolean leave(int ll, int rr, int l, int r) {
        return ll > r || rr < l;
    }

    public void update(int ll, int rr, int l, int r, FlaggedCommutativeUndoOperation item) {
        if (leave(ll, rr, l, r)) {
            return;
        }
        if (enter(ll, rr, l, r)) {
            modify(item);
            return;
        }
        pushDown();
        int m = DigitUtils.floorAverage(l, r);
        left.update(ll, rr, l, m, item);
        right.update(ll, rr, m + 1, r, item);
        pushUp();
    }

    public void solve(int l, int r, Callback consumer) {
        for (FlaggedCommutativeUndoOperation item : items) {
            item.apply();
        }
        if (l < r) {
            int mid = DigitUtils.floorAverage(l, r);
            left.solve(l, mid, consumer);
            right.solve(mid + 1, r, consumer);
        } else {
            consumer.accept(l);
        }

        Collections.reverse(items);
        for (FlaggedCommutativeUndoOperation item : items) {
            item.undo();
        }
    }

}
