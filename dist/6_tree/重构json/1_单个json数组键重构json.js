"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const foo = ['a', 'b', 'c', 1];
const gen = (arr) => {
    const head = arr.shift();
    return {
        [head]: arr.length <= 1 ? arr[0] : gen(arr),
    };
};
console.dir(gen(foo), { depth: null });
