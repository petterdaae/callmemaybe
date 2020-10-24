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
mov rax, [rsp+8]
push rax
mov rax, [rsp+8]
pop rbx
add rax, rbx
push rax
mov rax, 0x3
pop rbx
imul rax, rbx
pop rbx
pop rbx




mov rdi, format
mov rsi, rax
xor rax, rax
call printf

mov rax, 60
syscall

