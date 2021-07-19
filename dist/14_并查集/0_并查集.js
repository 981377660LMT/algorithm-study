"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UnionFind = void 0;
// 并查集是一种数据结构，并(Union)代表合并，查(Find)代表查找，用来查找两根元素是否具有公共的根。
// 集代表这是一个以字典为基础的数据结构，基本功能是合并集合中的元素，
// 查找集合中的元素。
// 并查集的典型应用是有关连通分量的问题，
// 并查集解决单个问题（添加，合并，查找）的时间复杂度都是O(1)(因为都是用的map的set和get方法)。
class UnionFind {
    constructor() {
        this.parent = new Map();
        this.count = 0;
    }
    /**
     *
     * @param key 把一个新节点添加到并查集中，它的父节点应该为UnionFind.rootSymbol。
     */
    add(key) {
        if (!this.parent.has(key)) {
            this.parent.set(key, UnionFind.rootSymbol);
            this.count++;
        }
        return this;
    }
    /**
     *
     * @description 如果两个节点是连通的，那么就要把他们合并，也就是他们的祖先是相同的。
     * @example
     * ```js
     * const union = new UnionFind<number>()
     * union.add(1).add(2).add(3).add(4).union(2, 3).union(4, 3)
     * console.dir(union, { depth: null })
     *
     * // output:
     * UnionFind {
     *   parent: Map(4) { 1 => undefined, 2 => 3, 3 => undefined, 4 => 3 }
     * }
     * ```
     */
    union(key1, key2) {
        const root1 = this.find(key1);
        const root2 = this.find(key2);
        if (root1 !== undefined && root2 !== undefined && root1 !== root2) {
            this.parent.set(root1, root2);
            this.count--;
        }
        return this;
    }
    /**
     * @description 判断两个节点是否处于同一个连通分量的时候，就要判断他们的祖先是否相同。
     */
    isConnected(key1, key2) {
        return this.find(key1) === this.find(key2);
    }
    /**
     *
     * @param key 查找祖先；如果节点的父节点不为空或者symbol，那就不断迭代。
     * @returns 返回undefined代表key不在并查集中
     */
    find(key) {
        let root = key;
        if (!this.parent.has(root))
            return undefined;
        while (this.parent.get(root) !== UnionFind.rootSymbol) {
            root = this.parent.get(root);
        }
        return root;
    }
}
exports.UnionFind = UnionFind;
UnionFind.rootSymbol = Symbol.for('UnionFind_Root');
if (require.main === module) {
    const union = new UnionFind();
    union.add(1).add(2).add(3).add(4).union(2, 3).union(4, 3).add(6);
    console.dir(union, { depth: null });
}
