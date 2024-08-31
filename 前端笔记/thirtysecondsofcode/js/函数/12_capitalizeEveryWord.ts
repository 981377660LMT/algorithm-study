const capitalizeEveryWord = (str: string) => str.replace(/\b[a-z]/g, (char: string) => char.toUpperCase())
capitalizeEveryWord('hello world!') // 'Hello World!'
