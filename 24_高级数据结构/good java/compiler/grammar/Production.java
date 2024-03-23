package template.compiler.grammar;

import template.compiler.ast.AstNode;

import java.util.Iterator;

public abstract class Production implements Iterable<Symbol> {
    Symbol[] symbols;

    public Production(Symbol... symbols) {
        this.symbols = symbols;
    }

    @Override
    public Iterator<Symbol> iterator() {
        return new Iterator<Symbol>() {
            int index = 0;

            @Override
            public boolean hasNext() {
                return index < symbols.length;
            }

            @Override
            public Symbol next() {
                return symbols[index++];
            }
        };
    }

    public int length() {
        return symbols.length;
    }

    public Symbol get(int i) {
        return symbols[i];
    }

    public abstract AstNode parse(AstNode... children);
}
