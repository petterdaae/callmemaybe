extern printf
extern malloc
extern free

global main

section .data
format: db '%d', 10, 0

section .text

main:
push rbx


mov rdi, 64
call malloc


mov rbx, 72
mov [rax], rbx

mov rbx, [rax]
mov rax, rbx



mov rdi, format
mov rsi, rax
xor rax, rax
call printf



pop rbx
mov rax, 0
ret

