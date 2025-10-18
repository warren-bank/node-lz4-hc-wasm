import fs   from 'node:fs'
import path from 'node:path'
import {fileURLToPath} from 'node:url'

let lz4 = null

export const init = async () => {
  if (lz4) return lz4

  if (!globalThis.fs)      globalThis.fs      = fs
  if (!globalThis.path)    globalThis.path    = path
  if (!globalThis.process) globalThis.process = process

  const __filename = fileURLToPath(import.meta.url)
  const __dirname  = path.dirname(__filename)

  const wasmPath  = path.resolve(__dirname, '../lz4.wasm')
  const wasmBytes = fs.readFileSync(wasmPath)

  await import('../wasm_exec.js')

  const go = new globalThis.Go()

  const {instance} = await WebAssembly.instantiate(wasmBytes, go.importObject)
  go.run(instance)

  const compressBlockBound      = globalThis.compressBlockBound
  const uncompressBlock         = globalThis.uncompressBlock
  const uncompressBlockWithDict = globalThis.uncompressBlockWithDict
  const compressBlock           = globalThis.compressBlock
  const compressBlockHC         = globalThis.compressBlockHC

  delete globalThis.Go
  delete globalThis.compressBlockBound
  delete globalThis.uncompressBlock
  delete globalThis.uncompressBlockWithDict
  delete globalThis.compressBlock
  delete globalThis.compressBlockHC

  lz4 = {
    compressBlockBound,
    uncompressBlock,
    uncompressBlockWithDict,
    compressBlock,
    compressBlockHC
  }

  return lz4
}
