"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const sets = [[0], [0, 4], [1, 3], [2], [3, 4]];
class UnionFind {
    constructor() {
        this.map = new Map();
    }
    union(key1, key2) {
        const root1 = this.find(key1);
        const root2 = this.find(key2);
        if (root1 !== root2) {
            this.map.set(root1, root2);
        }
    }
    // 不存在则直接返回key
    find(key) {
        while (this.map.has(key)) {
            key = this.map.get(key);
        }
        return key;
    }
}
const unionFind = new UnionFind();
for (let i = 0; i < sets.length; i++) {
    for (let j = 0; j < sets[i].length - 1; j++) {
        unionFind.union(sets[i][j], sets[i][j + 1]);
    }
}
console.log(unionFind);
console.log(unionFind.find(0));
console.log(unionFind.find(1));
console.log(unionFind.find(2));
console.log(unionFind.find(3));
console.log(unionFind.find(4));
