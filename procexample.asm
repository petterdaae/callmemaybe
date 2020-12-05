; A = 42
; proc1 = <> => {
;     println A + 1
; }
; B = 37
; call proc1
; C= 12
; proc2 = <> => {
;     D = 8
;     println B + D
; }
; call proc1
; call proc2
; E = 14
; call proc2

extern printf

global main

section .data
format: db '%d', 10, 0

section .text

main:
    ; Make sure rsp is initialized
    push rbx

    ; A
    mov rax, 42
    push rax

    ; Init proc1 here, store context
    mov rax, [rsp+0]
    push rax
    ; Store stack size in compiler = 2

    ; B
    mov rax, 37 
    push rax

    ; Move current stack size to rcx = 3
    mov rcx, 3
    call proc1

    ; C
    mov rax, 12
    push rax

    ; Init proc2 here, store context
    mov rax, [rsp+0]
    push rax
    mov rax, [rsp+8]
    push rax
    mov rax, [rsp+16]
    push rax

    ; Move current stack size to rcx = 4
    mov rcx, 7
    call proc1

    mov rcx, 7
    call proc2

    ; E
    mov rax, 14
    push rax

    mov rcx, 8
    call proc2
    
    ; Stack should be empty at end of program
    pop rax
    pop rax
    pop rax
    pop rbx
    pop rax
    pop rax
    pop rax
    pop rax

    mov rax, 0
    ret

proc1:
    ; Sub current stack size with init stack size
    mov rdx, 2
    sub rcx, rdx
    imul rcx, 8


    ; Add one to A
    mov rax, 1
    mov rbx, [rsp+rcx+0+8] ; Add one more (8) because of return pointer
    add rax, rbx

    ; Print rax
    call printrax

    ret

proc2:
    ;Sub current stack size with init stack size
    mov rdx, 4
    sub rcx, rdx
    imul rcx, 8

    ; D
    mov rax, 8
    push rax

    ; Add D to B
    mov rax, [rsp+rcx+8+8+8] ; [rsp + rcx + relative to init + return pointer + one more element on stack in this proc]
    mov rbx, [rsp+0] ; Keep track of variables local to the procedure
    add rax, rbx

    ; Print rax
    call printrax

    ; Stack should be the same size as start of procedure
    pop rax

    ret
    
printrax:
    push rcx
    push rdx
    push rdi
    push rsi
    mov rdi, format
    mov rsi, rax
    xor rax, rax
    call printf
    pop rsi
    pop rdi
    pop rdx
    pop rcx
    ret

