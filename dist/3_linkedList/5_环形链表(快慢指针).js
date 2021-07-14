"use strict";
// 带环:两个人同时起跑，速度快的会超过速度慢的;
// 不带环:遍历结束后，没有相逢
Object.defineProperty(exports, "__esModule", { value: true });
class Node {
    constructor(value, next) {
        this.value = value;
        this.next = next;
    }
}
const a = new Node(1);
const b = new Node(1);
const c = new Node(2);
const d = new Node(3);
a.next = b;
b.next = c;
c.next = d;
d.next = b;
const hasCycle = (node) => {
    let fastNode = node;
    let slowNode = node;
    while (fastNode) {
        fastNode = fastNode.next?.next;
        slowNode = slowNode?.next;
        if (slowNode === fastNode) {
            return true;
        }
    }
    return false;
};
console.log(hasCycle(a));
console.dir(a, { depth: null });
