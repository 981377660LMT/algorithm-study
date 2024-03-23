package template.string.re;

import java.util.Collection;

public interface State {
    Transfer next(int c);

    Collection<Transfer> adj();

    void register(Transfer s);

    int id();

    static boolean isInvalid(State s) {
        return s.id() == 0;
    }

    static boolean isAccept(State s) {
        return s.id() == 1;
    }
}
