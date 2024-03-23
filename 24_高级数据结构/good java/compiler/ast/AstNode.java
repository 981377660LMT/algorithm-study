package template.compiler.ast;

import template.compiler.context.Context;

public interface AstNode<T> {
    T eval(Context ctx);
}
