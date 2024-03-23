package template.string.re;

import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Deque;
import java.util.List;

public class Parser {
    static boolean[] forbiddenNothing = new boolean[128];

    public static NFARegularExpression parse(CharSequence cs) {
        return parse(cs, forbiddenNothing);
    }

    public static NFARegularExpression parse(CharSequence cs, boolean[] forbidden) {
        assert forbidden.length >= 128;
        return new ParserHelper(cs, forbidden).parse();
    }

    /**
     * special character:
     * .
     * *
     * +
     * ?
     * |
     * (
     * )
     * \
     */
    private static class ParserHelper {
        CharSequence str;
        int rpos = 0;
        State accept;
        State invalid;
        int n;
        Subgraph or = new Subgraph(null, null);
        Subgraph bracket = new Subgraph(null, null);
        List<State> all;
        boolean[] forbidden;
        public ParserHelper(CharSequence str, boolean[] forbidden) {
            this.forbidden = forbidden;
            this.str = str;
            this.n = str.length();
            invalid = new InvalidState();
            accept = new AcceptState(new TransferImpl(invalid));
            all = new ArrayList<>(n + 2);
            all.add(invalid);
            all.add(accept);
        }

        private int read() {
            if (rpos >= n) {
                return -1;
            }
            return str.charAt(rpos++);
        }

        private void throwExceptionBecauseOfLastCharacter() {
            throw new IllegalArgumentException("Invalid regular expression because of the " + rpos + "-th character");
        }

        private void nextId(AbstractState state) {
            state.setId(all.size());
            all.add(state);
        }

        Deque<Subgraph> dq;

        private Subgraph eval(Subgraph s) {
            if (s == null) {
                AbstractState state = new MatchNothing(new TransferImpl(invalid));
                nextId(state);
                Transfer t = new TransferImpl();
                state.register(t);
                s = new Subgraph(state, t);
            }
            return s;
        }

        private Subgraph and(Subgraph a, Subgraph b) {
            if (a == null) {
                return b;
            }
            if (b == null) {
                return a;
            }
            a.outbound.set(b.head);
            return new Subgraph(a.head, b.outbound);
        }

        private Subgraph or(Subgraph a, Subgraph b) {
            a = eval(a);
            if (b == null) {
                return a;
            }
            AbstractState branch = new MatchNothing(new TransferImpl(invalid));
            nextId(branch);
            Transfer t1 = new TransferImpl(a.head);
            Transfer t2 = new TransferImpl(b.head);
            branch.register(t1);
            branch.register(t2);
            return new Subgraph(branch, new MultiTransfer(a.outbound, b.outbound));
        }

        private void compressTail(boolean removeBracket) {
            Subgraph lastOr = null;
            Subgraph lastAnd = null;

            while (!dq.isEmpty()) {
                Subgraph tail = dq.removeLast();
                if (tail == or || tail == bracket) {
                    lastOr = or(lastAnd, lastOr);
                    lastAnd = null;

                    if (tail == bracket) {
                        if (!removeBracket) {
                            dq.addLast(tail);
                        }
                        removeBracket = false;
                        break;
                    }
                    continue;
                }
                lastAnd = and(tail, lastAnd);
            }
            if (removeBracket) {
                throwExceptionBecauseOfLastCharacter();
            }
            dq.addLast(eval(lastOr));
        }

        private Subgraph assertRemoveLast() {
            if (dq.isEmpty()) {
                throwExceptionBecauseOfLastCharacter();
            }
            Subgraph tail = dq.removeLast();
            if (tail == or || tail == bracket) {
                throwExceptionBecauseOfLastCharacter();
            }
            return tail;
        }

        public NFARegularExpression parse() {
            dq = new ArrayDeque<>(n + 2);
            dq.addLast(bracket);
            int nc;
            boolean escape = false;
            while ((nc = read()) >= 0) {
                Subgraph next = null;
                if (escape) {
                    escape = false;
                    Transfer outbound = new TransferImpl();
                    AbstractState state = new CharacterState(new TransferImpl(invalid), outbound, nc);
                    nextId(state);
                    next = new Subgraph(state, outbound);
                } else if (nc == '\\' && !forbidden[nc]) {
                    escape = true;
                    continue;
                } else if (nc == '.' && !forbidden[nc]) {
                    Transfer outbound = new TransferImpl();
                    AbstractState state = new MatchAnyState(outbound);
                    nextId(state);
                    next = new Subgraph(state, outbound);
                } else if (nc == '|' && !forbidden[nc]) {
                    //cool
                    compressTail(false);
                    dq.addLast(or);
                    continue;
                } else if (nc == '(' && !forbidden[nc]) {
                    dq.addLast(bracket);
                    continue;
                } else if (nc == ')' && !forbidden[nc]) {
                    compressTail(true);
                    continue;
                } else if (nc == '*' && !forbidden[nc]) {
                    Subgraph tail = assertRemoveLast();
                    AbstractState branch = new MatchNothing(new TransferImpl(invalid));
                    nextId(branch);
                    Transfer t1 = new TransferImpl();
                    Transfer t2 = new TransferImpl();
                    branch.register(t1);
                    branch.register(t2);
                    t1.set(tail.head);
                    tail.outbound.set(branch);
                    next = new Subgraph(branch, t2);
                } else if (nc == '+' && !forbidden[nc]) {
                    Subgraph tail = assertRemoveLast();
                    AbstractState branch = new MatchNothing(new TransferImpl(invalid));
                    nextId(branch);
                    Transfer t1 = new TransferImpl();
                    Transfer t2 = new TransferImpl();
                    branch.register(t1);
                    branch.register(t2);
                    t1.set(tail.head);
                    tail.outbound.set(branch);
                    next = new Subgraph(tail.head, t2);
                } else if (nc == '?' && !forbidden[nc]) {
                    Subgraph tail = assertRemoveLast();
                    AbstractState branch = new MatchNothing(new TransferImpl(invalid));
                    nextId(branch);
                    Transfer t1 = new TransferImpl();
                    Transfer t2 = new TransferImpl();
                    branch.register(t1);
                    branch.register(t2);
                    t1.set(tail.head);
                    next = new Subgraph(branch, new MultiTransfer(tail.outbound, t2));
                } else {
                    Transfer outbound = new TransferImpl();
                    AbstractState state = new CharacterState(new TransferImpl(invalid), outbound, nc);
                    nextId(state);
                    next = new Subgraph(state, outbound);
                }

                dq.addLast(next);
            }

            compressTail(true);
            if (dq.size() != 1) {
                throwExceptionBecauseOfLastCharacter();
            }
            if (escape) {
                throwExceptionBecauseOfLastCharacter();
            }

            Subgraph g = assertRemoveLast();
            g.outbound.set(accept);
            return new NFARegularExpression(all.toArray(new State[0]), g.head);
        }
    }
}
