extern printf

global main

section .text
main:
    mov rdi, format
    mov rsi, 0x42
    xor rax, rax
    call printf
    mov rax, 60
    syscall

section .data
    format: db '%x',10,0

