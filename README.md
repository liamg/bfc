# bfc

*bfc* is a bad (probably) [Brainfuck](https://en.wikipedia.org/wiki/Brainfuck) compiler. 

It compiles only to *x64* assembly, though there are flags that you can use to get your code assembled and linked into a standard Linux elf executable.

I wrote this compiler because I had no idea how to write compilers, and thought it might be interesting. It was very interesting and I still have no idea how to write compilers because Brainfuck is a very syntactically simple language and has very few possible operations. I'll be working on a C compiler next instead. I didn't think this through.

## Installation

Grab the binary from the [latest release](https://github.com/liamg/bfc/releases/latest).

## Usage

### Compile Brainfuck to x64 asm:

```bash
bfc -o hello.asm hello.bf
```

### Compile and assemble Brainfuck to x64 object:

```bash
bfc -a -o hello.o hello.bf
```

### Compile, assemble and link Brainfuck to ELF binary:

```bash
bfc -al -o hello hello.bf
```

### All of the above and run the resultant binary

```bash
bfc -r hello.bf
```
