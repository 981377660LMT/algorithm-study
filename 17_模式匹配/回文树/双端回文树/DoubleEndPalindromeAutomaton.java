package template.string;

import java.util.ArrayList;
import java.util.List;
import java.util.function.Consumer;

public class DoubleEndPalindromeAutomaton {
    final int minCharacter;
    final int maxCharacter;
    int range;

    Node odd;
    Node even;

    char[] data;
    int zero;
    int frontSize;
    int backSize;
    Node frontBuildLast;
    Node backBuildLast;
    List<Node> all;

    long palindromeSubstringCnt = 0;

    private Node newNode() {
        Node ans = new Node(range);
        all.add(ans);
        return ans;
    }

    public DoubleEndPalindromeAutomaton(int minCharacter, int maxCharacter, int frontAddition, int backAddition) {
        int cap = frontAddition + backAddition;
        all = new ArrayList<>(2 + cap);
        data = new char[cap];
        zero = frontAddition;
        frontSize = zero - 1;
        backSize = zero;

        odd = newNode();
        odd.len = -1;
        odd.depth = 0;

        even = newNode();
        even.fail = odd;
        even.len = 0;
        even.depth = 0;

        all.clear();
        backBuildLast = frontBuildLast = odd;
    }

    public void buildFront(char c) {
        data[frontSize--] = c;

        int index = c - minCharacter;

        Node trace = frontBuildLast;
        while (frontSize + 2 + trace.len >= backSize) {
            trace = trace.fail;
        }

        while (data[frontSize + trace.len + 2] != c) {
            trace = trace.fail;
        }

        if (trace.next[index] != null) {
            frontBuildLast = trace.next[index];
        } else {
            Node now = newNode();
            now.len = trace.len + 2;
            trace.next[index] = now;

            if (now.len == 1) {
                now.fail = even;
            } else {
                trace = trace.fail;
                while (data[frontSize + trace.len + 2] != c) {
                    trace = trace.fail;
                }
                now.fail = trace.next[index];
            }
            now.depth = now.fail.depth + 1;
            frontBuildLast = now;
        }
        if(frontBuildLast.len == backSize - frontSize - 1){
            backBuildLast = frontBuildLast;
        }
        palindromeSubstringCnt += frontBuildLast.depth;
    }

    public void buildBack(char c) {
        data[backSize++] = c;

        int index = c - minCharacter;

        Node trace = backBuildLast;
        while (backSize - 2 - trace.len <= frontSize) {
            trace = trace.fail;
        }

        while (data[backSize - trace.len - 2] != c) {
            trace = trace.fail;
        }

        if (trace.next[index] != null) {
            backBuildLast = trace.next[index];
        } else {
            Node now = newNode();
            now.len = trace.len + 2;
            trace.next[index] = now;

            if (now.len == 1) {
                now.fail = even;
            } else {
                trace = trace.fail;
                while (data[backSize - trace.len - 2] != c) {
                    trace = trace.fail;
                }
                now.fail = trace.next[index];
            }
            now.depth = now.fail.depth + 1;
            backBuildLast = now;
        }
        if(backBuildLast.len == backSize - frontSize - 1){
            frontBuildLast = backBuildLast;
        }
        palindromeSubstringCnt += backBuildLast.depth;
    }

    public void visit(Consumer<Node> consumer) {
        for (int i = all.size() - 1; i >= 0; i--) {
            consumer.accept(all.get(i));
        }
    }

    public long palindromeSubstringCnt() {
        return palindromeSubstringCnt;
    }

    public int distinctPalindromeSubstring() {
        return all.size();
    }

    public static class Node {
        public Node(int range) {
            next = new Node[range];
        }
        public Node[] next;
        public Node fail;
        public int len;
        public int depth;
    }
}
