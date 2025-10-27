import {init} from '../../dist/esm/lz4.mjs'

import fs   from 'node:fs'
import path from 'node:path'
import {fileURLToPath} from 'node:url'

const run_tests = async () => {
  const lz4 = await init()

  const __filename = fileURLToPath(import.meta.url)
  const __dirname  = path.dirname(__filename)

  const json = fs.readFileSync(
    path.resolve(__dirname, '../../package.json'),
    {encoding: 'utf8'}
  )

  const data = {lz4, json}

  await test_01(data); print_test_sep()
  await test_02(data); print_test_sep()
  await test_03(data); print_test_sep()
  await test_04(data, 0); print_test_sep()
  await test_04(data, 1); print_test_sep()
  await test_04(data, 2); print_test_sep()
  await test_04(data, 3); print_test_sep()
  await test_04(data, 4); print_test_sep()
  await test_04(data, 5); print_test_sep()
  await test_04(data, 6); print_test_sep()
  await test_04(data, 7); print_test_sep()
  await test_04(data, 8); print_test_sep()
  await test_04(data, 9); print_test_sep()
}

const print_test_sep = () => {
  console.log('-'.repeat(40))
}

const test_01 = async (data) => {
  const {lz4, json} = data

  console.log('lz4.uncompressBlock =',         typeof lz4.uncompressBlock)
  console.log('lz4.uncompressBlockWithDict =', typeof lz4.uncompressBlockWithDict)
  console.log('lz4.compressBlock =',           typeof lz4.compressBlock)
  console.log('lz4.compressBlockHC =',         typeof lz4.compressBlockHC)
}

const test_02 = async (data) => {
  const {lz4, json} = data

  console.log('json =', json.trim())
}

const test_03 = async (data) => {
  const {lz4, json} = data

  const src_1 = (new TextEncoder()).encode(json)
  const dst_1 = await lz4.compressBlock(src_1)
  console.log('lz4.compressBlock =>', dst_1.constructor.name, 'with', dst_1.byteLength, 'bytes')

  const dst_2 = await lz4.uncompressBlock(dst_1)
  console.log('lz4.uncompressBlock =>', dst_2.constructor.name, 'with', dst_2.byteLength, 'bytes')

  const src_2 = (new TextDecoder('utf-8')).decode(dst_2)
  console.log(`after compression and decompression, json is ${(json === src_2) ? '' : 'not '}the same`)
}

const test_04 = async (data, compressionLevel) => {
  const {lz4, json} = data

  const src_1 = (new TextEncoder()).encode(json)
  const dst_1 = await lz4.compressBlockHC(src_1, compressionLevel)
  console.log(`lz4.compressBlockHC(compressionLevel=${compressionLevel}) =>`, dst_1.constructor.name, 'with', dst_1.byteLength, 'bytes')

  const dst_2 = await lz4.uncompressBlock(dst_1)
  console.log('lz4.uncompressBlock =>', dst_2.constructor.name, 'with', dst_2.byteLength, 'bytes')

  const src_2 = (new TextDecoder('utf-8')).decode(dst_2)
  console.log(`after compression and decompression, json is ${(json === src_2) ? '' : 'not '}the same`)
}

run_tests()
