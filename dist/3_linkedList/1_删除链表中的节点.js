"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Node = void 0;
class Node {
    constructor(value, next) {
        this.value = value;
        this.next = next;
    }
}
exports.Node = Node;
const a = new Node(1);
const b = new Node(2);
const c = new Node(3);
a.next = b;
b.next = c;
// 无法获取被删除节点的上个节点
// 将被删除节点转移到下个节点
const deleteNode = (node) => {
    node.value = node.next?.value;
    node.next = node.next?.next;
};
console.dir(a);
deleteNode(b);
console.log(a);
// O(1)
