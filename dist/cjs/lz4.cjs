const fs   = require('fs')
const path = require('path')

let lz4 = null

const init = async () => {
  if (lz4) return lz4

  if (!globalThis.fs)      globalThis.fs      = fs
  if (!globalThis.path)    globalThis.path    = path
  if (!globalThis.process) globalThis.process = process

  const wasmPath  = path.join(__dirname, '../lz4.wasm')
  const wasmBytes = fs.readFileSync(wasmPath)

  require('../wasm_exec.js')

  const go = new globalThis.Go()

  const {instance} = await WebAssembly.instantiate(wasmBytes, go.importObject)
  go.run(instance)

  const compressFrame           = globalThis.compressFrame
  const uncompressFrame         = globalThis.uncompressFrame
  const compressBlockBound      = globalThis.compressBlockBound
  const compressBlock           = globalThis.compressBlock
  const compressBlockHC         = globalThis.compressBlockHC
  const uncompressBlock         = globalThis.uncompressBlock
  const uncompressBlockWithDict = globalThis.uncompressBlockWithDict

  delete globalThis.Go
  delete globalThis.compressFrame
  delete globalThis.uncompressFrame
  delete globalThis.compressBlockBound
  delete globalThis.compressBlock
  delete globalThis.compressBlockHC
  delete globalThis.uncompressBlock
  delete globalThis.uncompressBlockWithDict

  lz4 = {
    compressFrame,
    uncompressFrame,
    compressBlockBound,
    compressBlock,
    compressBlockHC,
    uncompressBlock,
    uncompressBlockWithDict
  }

  return lz4
}

module.exports = {init}
