package template.string.re;

import java.util.Collection;
import java.util.Collections;

public class AcceptState implements State {
    Transfer invalid;

    public AcceptState(Transfer invalid) {
        this.invalid = invalid;
    }

    @Override
    public Transfer next(int c) {
        return invalid;
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
        return 1;
    }

    @Override
    public String toString() {
        return "";
    }
}
