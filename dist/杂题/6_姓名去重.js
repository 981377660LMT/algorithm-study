"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const nameDeduplicate = (names) => {
    if (names.length <= 1)
        return names;
    return [...new Set(names.map(name => name.toLowerCase()))];
};
console.log(nameDeduplicate([
    'James',
    'james',
    'Bill Gates',
    'bill Gates',
    'Hello World',
    'HELLO WORLD',
    'Helloworld',
]));
