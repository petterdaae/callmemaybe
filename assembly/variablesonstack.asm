; a = 1
; b = 2
; c = 3
; println b
extern printf

global main

section .data
    format: db '%x', 10, 0

section .text

main:
    mov rax, 0x1
    push rax
    mov rax, 0x2
    push rax
    mov rax, 0x3
    push rax
    mov rax, [rsp+8]
    mov rdi, format
    mov rsi, rax
    call printf
    mov rax, 60
    syscall

