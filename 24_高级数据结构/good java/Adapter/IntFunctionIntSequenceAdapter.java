package template.string;

import template.primitve.generated.datastructure.IntToIntegerFunction;

public class IntFunctionIntSequenceAdapter implements IntSequence {
    IntToIntegerFunction data;
    int l;
    int r;

    public IntFunctionIntSequenceAdapter(IntToIntegerFunction data, int l, int r) {
        this.data = data;
        this.l = l;
        this.r = r;
    }

    @Override
    public int get(int i) {
        return data.apply(l + i);
    }

    @Override
    public int length() {
        return r - l + 1;
    }

    @Override
    public IntSequence subsequence(int l, int r) {
        return new IntFunctionIntSequenceAdapter(data, this.l + l, this.l + r);
    }
}
