"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const foo = ['a', 'b', 'c', 'x', 'y', 1];
const bar = ['a', 'b', 'c', 'k', 'g', 2];
const e = ['e', 'f', 3];
const resJson = {};
const gen2 = (arr) => {
    const head = arr.shift();
    console.log(777, arr, head);
    return {
        [head]: arr.length <= 1 ? arr[0] : gen2(arr),
    };
};
const gen = (arr, p) => {
    if (!arr.length)
        return;
    const head = arr.shift();
    console.log(666, p[head], head, p, arr);
    p[head] = head in p ? { ...p[head], ...gen(arr, p[head]) } : gen2(arr);
};
gen(foo, resJson);
gen(bar, resJson);
// gen(e, resJson)
// console.dir(gen(foo, resJson), { depth: null })
console.dir(resJson, { depth: null });
console.log(resJson['a']['b']['c']['k']);
