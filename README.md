### [lz4-hc-wasm](https://github.com/warren-bank/node-lz4-hc-wasm)

LZ4 block format encoder/decoder: a WebAssembly implementation with variable high compression.

#### Credits

[Pierre Curto](https://github.com/pierrec) wrote [lz4](https://github.com/pierrec/lz4) as a Golang module.

This repo includes:
* his Golang module as a git submodule
* a [`main.go`](./src/main.go) wrapper to define the JavaScript API, and manage data type conversion
* a [build script](./bin/bash/2-build.sh) to generate a WASM library
* [dist](./dist)
  - the [WASM library](./dist/lz4.wasm)
  - a helper script to load the WASM library into JavaScript for each of the following environments:
    * [CommonJS](./dist/cjs/lz4.cjs)
    * [ESM](./dist/esm/lz4.mjs)
    * [web](./dist/web/lz4.js)
* a [test](./tests) for each environment

#### Installation

npm:
```bash
  npm install "@warren-bank/lz4-hc-wasm"
```

#### Usage

CommonJS:
```js
  const {init} = require('@warren-bank/lz4-hc-wasm')
```

ESM:
```js
  import {init} from '@warren-bank/lz4-hc-wasm'
```

common:
```js
  const lz4 = await init()
```

#### JavaScript API

```js
  type lz4.uncompressBlock         = (src: Uint8Array) => Promise<Uint8Array>
  type lz4.uncompressBlockWithDict = (src: Uint8Array, dict: Uint8Array) => Promise<Uint8Array>
  type lz4.compressBlock           = (src: Uint8Array) => Promise<Uint8Array>
  type lz4.compressBlockHC         = (src: Uint8Array, depth?: int) => Promise<Uint8Array>
```

where:
* `Promise<Uint8Array>` means that a Promise is returned, which will either:
  - resolve to an `Uint8Array`
  - reject if any error occurs

notes:
* test [results](./tests/cjs/test.log) confirm [this open issue](https://github.com/pierrec/lz4/issues/219),<br>which is that the `depth` compression-level parameter to `compressBlockHC` has (almost) no effect

#### To Do

* reduce the size of the WASM library by using a different Golang compiler
  - ex: [TinyGo](https://tinygo.org/)

#### Legal

* copyright: [Pierre Curto](https://github.com/pierrec)
* license: [BSD-3-Clause](https://github.com/pierrec/lz4/raw/cd9d7a4f66405a92f4881933816b828b0f7d5fe2/LICENSE)
