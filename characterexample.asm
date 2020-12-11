extern printf

global main

section .data

format: db '%c', 0

section .text

main:
    push rbx


    
    mov rax, 10
    mov rdi, format
    mov rsi, rax
    xor rax, rax
    call printf

    pop rbx
    mov rax, 0
    ret
