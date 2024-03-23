package template.string.re;

import java.util.Arrays;
import java.util.Collection;

public class MultiTransfer implements Transfer {
    Collection<Transfer> ts;

    public MultiTransfer(Collection<Transfer> ts) {
        this.ts = ts;
    }

    public MultiTransfer(Transfer... ts) {
        this(Arrays.asList(ts));
    }

    @Override
    public State get() {
        return ts.iterator().next().get();
    }

    @Override
    public void set(State state) {
        for (Transfer t : ts) {
            t.set(state);
        }
    }
}
