# CommentPurger

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

A simple and efficient command-line tool written in Go to remove comments from your code files.

CommentPurger helps you clean up your codebase for production builds or distribution by safely stripping comments from various file types. It recursively processes files in the directories you specify.

## Key Features

- **Multi-Language Support**: Works out of the box with HTML, CSS, JavaScript, TypeScript, and Vue (`.vue`) files.
- **Safe for HTML**: Uses a proper HTML parser to safely remove comments without breaking your markup.
- **Vue SFC Aware**: Intelligently parses `.vue` files, removing comments from `<template>`, `<script>`, and `<style>` sections correctly.
- **Recursive Processing**: You can point it at a single file or an entire directory, and it will find and process all supported files.
- **Fast**: Built with Go for high performance.

## Installation

To install `CommentPurger`, you need to have Go installed (version 1.18 or newer).

You can install the program directly using `go install`:

```bash
go install github.com/scamgi/commentpurger@latest
```

This will compile the program and place the `commentpurger` executable in your Go bin directory (`$GOPATH/bin`), which should be in your system's `PATH`.

## Usage

You can run the command against one or more files or directories.

#### To process a single file

```bash
commentpurger path/to/your/file.js
```

#### To process an entire directory recursively

```bash
commentpurger ./my-project-directory/
```

The program will modify the files in place.

## Building from Source

If you want to build the executable from the source code yourself:

1. **Clone the repository:**

    ```bash
    git clone https://github.com/scamgi/commentpurger.git
    cd commentpurger
    ```

2. **Build the program:**

    ```bash
    go build
    ```

    This will create the `commentpurger` executable in the current directory.

## Running Tests

To run the test suite and verify that everything is working as expected:

```bash
go test -v
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

