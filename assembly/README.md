# Assembly

- `sudo apt install as31 nasm`
- Assemble: `nasm -f elf64 hello.asm`
- Link: `gcc -no-pie -o out hello.o -lc`
- Run: `./out`

