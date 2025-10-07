require('fake-indexeddb/auto')
global.ResizeObserver = require('resize-observer-polyfill')
global.crypto ??= new (require('@peculiar/webcrypto').Crypto)()
global.FontFace = class FontFace {
	load() {
		return Promise.resolve()
	}
}

document.fonts = {
	add: () => {},
	delete: () => {},
	forEach: () => {},
	[Symbol.iterator]: () => [][Symbol.iterator](),
}

global.TextEncoder = require('util').TextEncoder
global.TextDecoder = require('util').TextDecoder

// Extract verson from package.json
const { version } = require('./package.json')

window.fetch = async (input, init) => {
	return {
		ok: true,
		json: async () => [],
	}
}
