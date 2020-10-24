extern printf
global main
section .data
format: db '%x', 10, 0
section .text
main:
mov rax, 0x3
push rax
mov rax, [rsp+0]
mov rdi, format
mov rsi, rax
xor rax, rax
call printf
mov rax, 0x154
push rax
mov rax, 0x2
push rax
mov rax, [rsp+0]
push rax
mov rax, 0x3
pop rbx
add rax, rbx
pop rbx
pop rbx
add rax, rbx
push rax
mov rax, [rsp+8]
push rax
mov rax, [rsp+8]
pop rbx
add rax, rbx
mov rdi, format
mov rsi, rax
xor rax, rax
call printf
mov rax, [rsp+8]
push rax
mov rax, [rsp+8]
pop rbx
imul rax, rbx
push rax
mov rax, [rsp+0]
mov rdi, format
mov rsi, rax
xor rax, rax
call printf
mov rax, 60
syscall

