extern printf
extern malloc
global main
section .data
format: db '%d', 10, 0
section .text

main:
    push rbx

    mov rdi, 24
    call malloc

    mov rdx, rax
    mov qword [rdx], 1
    mov qword [rdx+8], 2
    mov qword [rdx+8*2], 3
    mov qword [rdx+8*3], 4

    mov rax, [rdx]
    call printit
    mov rax, [rdx+8]
    call printit
    mov rax, [rdx+16]
    call printit


    pop rbx
    mov rax, 0
    ret

printit:
    push rdi
    push rsi
    push rax
    push rdx

    mov rdi, format
    mov rsi, rax
    xor rax, rax
    call printf
    
    pop rdx
    pop rax
    pop rsi
    pop rdi
    ret
