	extern printf
	extern malloc
	extern free
	global main
	section .date
	format: db '%d', 10, 0
	formatchar: db '%c', 10, 0
	formatcharnonewline: db '%c', 0
	section .text
	main:
	push rbx
	mov rdi, 16
	call malloc
	mov rdx, rax
	push rdx
	mov rdi, 8
	call malloc
	mov rdx, rax
	push rdx
	mov rax, 1
	pop rdx
	mov qword [rdx+0], rax
	mov rax, rdx
	pop rdx
	mov qword [rdx+0], rax
	push rdx
	mov rdi, 8
	call malloc
	mov rdx, rax
	push rdx
	mov rax, 3
	pop rdx
	mov qword [rdx+0], rax
	mov rax, rdx
	pop rdx
	mov qword [rdx+8], rax
	mov rax, rdx
	push rax
	mov rax, 0
	push rax
	jmp unique2
	unique1:
	mov rax, [rsp+0]
	push rax
	mov rax, [rsp+16]
	pop rcx
	mov rdx, rax
	mov rax, [rdx+8*rcx]
	push rax
	mov rax, 0
	push rax
	mov rax, [rsp+8]
	pop rcx
	mov rdx, rax
	mov rax, [rdx+8*rcx]
	push rax
	mov rax, [rsp+16]
	push rax
	mov rax, 1
	pop rbx
	add rax, rbx
	mov [rsp+16], rax
	pop rbx
	unique2:
	mov rax, [rsp+8]
	push rax
	mov rax, 1
	pop rbx
	cmp rbx, rax
	jl unique3
	jge unique4
	unique3:
	mov rax, 1
	jmp unique5
	unique4:
	mov rax, 0
	unique5:
	cmp rax, 1
	je unique1
	pop rbx
	pop rbx
	pop rbx
	pop rbx
	mov rax, 0
	ret

