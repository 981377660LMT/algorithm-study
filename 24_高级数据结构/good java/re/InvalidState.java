package template.string.re;

import java.util.Collection;
import java.util.Collections;

public class InvalidState implements State {
    Transfer self = new TransferImpl(this);

    @Override
    public Transfer next(int c) {
        return self;
    }

    @Override
    public Collection<Transfer> adj() {
        return Collections.emptyList();
    }

    @Override
    public void register(Transfer s) {
        throw new UnsupportedOperationException();
    }

    @Override
    public int id() {
        return 0;
    }

    @Override
    public String toString() {
        return "";
    }
}
