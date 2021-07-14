"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class Node {
    constructor(value, next) {
        this.value = value;
        this.next = next;
    }
}
const a = new Node(1);
const b = new Node(2);
const c = new Node(3);
a.next = b;
b.next = c;
const reverseList = (head) => {
    // 注意p1p2都是node
    let n1 = undefined;
    let n2 = head;
    while (n2) {
        //@ts-ignore
        // p2的下一个Node
        const tmp = n2.next;
        // 最重要的关系，为串联节点做准备
        n2.next = n1;
        n1 = n2;
        n2 = tmp;
    }
    return n1;
};
console.log(reverseList(a));
