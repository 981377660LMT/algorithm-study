"use strict";
class Stack {
    constructor() {
        this.queue = [];
        this.tmp = [];
    }
    push(x) {
        this.queue.push(x);
    }
    pop() {
        while (this.queue.length > 1) {
            this.tmp.push(this.queue.shift());
        }
        const ele = this.queue.shift();
        this.tmp.push(ele);
        this.queue = this.tmp;
        this.tmp = [];
        return ele;
    }
}
