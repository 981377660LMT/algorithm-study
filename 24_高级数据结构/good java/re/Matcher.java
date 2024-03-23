package template.string.re;

import java.util.BitSet;
import java.util.Optional;

public class Matcher {
    State head;
    State[] states;
    BitSet prev;
    BitSet next;
    int[] leftPrev;
    int[] leftNext;

    public Matcher(State[] states, State head) {
        this.states = states;
        this.head = head;
        prev = new BitSet(states.length);
        next = new BitSet(states.length);
    }

    void dfsUpdate(State root) {
        if (next.get(root.id())) {
            return;
        }
        next.set(root.id());
        for (Transfer t : root.adj()) {
            dfsUpdate(t.get());
        }
    }

    void dfsUpdate(State root, int v) {
        if (next.get(root.id())) {
            return;
        }
        next.set(root.id());
        leftNext[root.id()] = v;
        for (Transfer t : root.adj()) {
            dfsUpdate(t.get(), v);
        }
    }

    void swap() {
        {
            BitSet tmp = prev;
            prev = next;
            next = tmp;
        }
        {
            int[] tmp = leftNext;
            leftNext = leftPrev;
            leftPrev = tmp;
        }
    }

    /**
     * 判断s[start,end)是否匹配正则表达式，O(n(end-start))
     *
     * @param s
     * @return
     */
    public boolean match(CharSequence s, int start, int end) {
        next.clear();
        dfsUpdate(head);
        swap();
        for (int i = start; i < end; i++) {
            next.clear();
            char c = s.charAt(i);
            for (int j = prev.nextSetBit(0); j >= 0; j = prev.nextSetBit(j + 1)) {
                dfsUpdate(states[j].next(c).get());
            }
            swap();
        }
        return prev.get(1);
    }

    /**
     * 找到s[start,end)中第一个匹配正则表达式的子串s[a,b)，其中b最小，a任意
     * 如果找到时间复杂度为O(n(b-start))，否则时间复杂度为O(n(end-start))
     *
     * @param s
     * @return
     */
    public Optional<int[]> find(CharSequence s, int start, int end) {
        if(leftNext == null){
            leftPrev = new int[states.length];
            leftNext = new int[states.length];
        }
        next.clear();
        dfsUpdate(head, start);
        swap();
        for (int i = start; i <= end; i++) {
            if (prev.get(1)) {
                //find
                return Optional.of(new int[]{leftPrev[1], i});
            }
            if (i == end) {
                break;
            }
            next.clear();
            char c = s.charAt(i);
            for (int j = prev.nextSetBit(0); j >= 0; j = prev.nextSetBit(j + 1)) {
                dfsUpdate(states[j].next(c).get(), leftPrev[j]);
            }
            dfsUpdate(head, i + 1);
            swap();
        }
        return Optional.empty();
    }
}
