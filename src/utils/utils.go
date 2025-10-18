package utils

// https://pkg.go.dev/syscall/js

import (
  "syscall/js"
)

// -------------------------------------------------------------------

func ValidateJsUint8Array(val js.Value) (bool) {
  return (val.Type() == js.TypeObject) && val.InstanceOf(js.Global().Get("Uint8Array"));
}

func ReadJsBytes(jsBytes js.Value) ([]byte) {
  length := jsBytes.Get("length").Int()
  goBytes := make([]byte, length)
  js.CopyBytesToGo(goBytes, jsBytes)
  return goBytes
}

func WriteJsBytes(goBytes []byte) (js.Value) {
  jsBytes := js.Global().Get("Uint8Array").New(len(goBytes))
  js.CopyBytesToJS(jsBytes, goBytes)
  return jsBytes
}

// -------------------------------------------------------------------

func ValidateJsNumber(val js.Value) (bool) {
  return (val.Type() == js.TypeNumber);
}

func ReadJsInt(jsNumber js.Value) (int) {
  return jsNumber.Int()
}

func ReadJsUint32(jsNumber js.Value) (uint32) {
  return uint32(jsNumber.Int())
}

// -------------------------------------------------------------------

func ReturnJsUint8ArrayAsPromise(dst []byte, n int, err error) (js.Value) {
  if err != nil {
    return ReturnAsPromise(nil, err)
  } else {
    dst = dst[:n]
    jsBytes := WriteJsBytes(dst)
    return ReturnAsPromise(jsBytes, nil)
  }
}

func ReturnAsPromise(val any, err error) (js.Value) {
  handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    resolve := args[0]
    reject  := args[1]

    go func() {
      if err != nil {
        errorObject := js.Global().Get("Error").New(err.Error())
        reject.Invoke(errorObject)
      } else {
        resolve.Invoke(js.ValueOf(val))
      }
    }()

    return nil
  })

  return js.Global().Get("Promise").New(handler)
}
