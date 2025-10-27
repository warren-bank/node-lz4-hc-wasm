### [lz4-hc-wasm](https://github.com/warren-bank/wasm-lz4-hc)

WebAssembly implementation for LZ4 high compression encoder and decoder with support for both frame-level and block-level data formats.

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
  // frame-level data format
  type lz4.compressFrame           = (src: Uint8Array, depth?: int) => Promise<Uint8Array>
  type lz4.uncompressFrame         = (src: Uint8Array) => Promise<Uint8Array>

  // block-level data format
  type lz4.compressBlockBound      = (n: int) => int
  type lz4.compressBlock           = (src: Uint8Array, dstSize?: int) => Promise<Uint8Array>
  type lz4.compressBlockHC         = (src: Uint8Array, depth: int, dstSize?: int) => Promise<Uint8Array>
  type lz4.uncompressBlock         = (src: Uint8Array, dstSize?: int) => Promise<Uint8Array>
  type lz4.uncompressBlockWithDict = (src: Uint8Array, dict: Uint8Array, dstSize?: int) => Promise<Uint8Array>
```

where:

* [frame-level data format](https://github.com/lz4/lz4/blob/dev/doc/lz4_Frame_format.md) includes:
  - a header
    * magic number
    * descriptor fields
  - a sequence of blocks of compressed data
  - a footer
    * end mark
    * checksum
* block-level data format includes:
  - a single block of compressed data
* `Promise<Uint8Array>` means that a Promise is returned, which will either:
  - resolve to an `Uint8Array`
  - reject if any error occurs
* the following block-level details are low-level and can generally be ignored:
  - the function: `lz4.compressBlockBound`
  - the parameter: `dstSize?: int`
* the `depth` parameter is used to control the level of compression
  - its value can range: `0` to `9`
  - for frame-level compression, this value defaults to `0`
    * `lz4.compressBlock` is used when this value is equal to `0`
    * `lz4.compressBlockHC` is used when this value is greater than `0`
  - for block-level high compression, this is a required value
    * note that in this case, `lz4.compressBlock` is _not_ used when this value is equal to `0`

notes:
* test [results](./tests/cjs/test.log) confirm [this open issue](https://github.com/pierrec/lz4/issues/219),<br>which is that the `depth` compression-level parameter to `compressBlockHC` has no effect

#### To Do

* reduce the size of the WASM library by using a different Golang compiler
  - ex: [TinyGo](https://tinygo.org/)

#### Legal

* copyright: [Pierre Curto](https://github.com/pierrec)
* license: [BSD-3-Clause](https://github.com/pierrec/lz4/raw/cd9d7a4f66405a92f4881933816b828b0f7d5fe2/LICENSE)
