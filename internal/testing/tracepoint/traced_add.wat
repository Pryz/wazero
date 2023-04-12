(module
  (import "tracepoint" "call_add" (func $__imported_tracepoint_call_add (param i32 i32)))
  (import "tracepoint" "return_add" (func $__imported_tracepoint_return_add (param i32) (result i32)))

  (func $add (param i32) (param i32) (result i32)
    local.get 0
    local.get 1
    call $__imported_tracepoint_call_add
    local.get 0
    local.get 1
    i32.add
    call $__imported_tracepoint_return_add)
  (export "add" (func $add))
)
