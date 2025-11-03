# CommentPurger

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**CommentPurger** is a blazingly fast, command-line utility for stripping comments from your code. Built with Go, it's a simple and efficient tool designed to help you clean up your codebase for production builds, distribution, or analysis by safely removing comments from a variety of file types.

It recursively processes entire directories, modifying files in-place to prepare them for their next stage.

---

## Why Use CommentPurger?

-   **Production Builds**: Reduce file sizes by removing unnecessary comments before deploying your web application.
-   **Code Distribution**: Share your code without including internal notes or commented-out legacy code.
-   **Code Analysis**: Prepare source files for static analysis tools that might be affected by comments.

## Key Features

-   **High Performance**: Written in Go for maximum speed and efficiency.
-   **Recursive Processing**: Point it at a directory, and it will find and process all supported files within it.
-   **Language-Aware Parsing**: Uses appropriate methods for different file types to ensure comments are removed safely.
    -   **Safe for HTML**: Employs a proper HTML parser to avoid breaking your markup.
    -   **Vue SFC Aware**: Intelligently handles `.vue` Single-File Components, cleaning `<template>`, `<script>`, and `<style>` sections correctly.
-   **In-Place Modification**: The tool directly updates your files, making it seamless to integrate into a build script.

---

## Supported Languages

CommentPurger works out of the box with the most common web development file types:

-   HTML (`.html`)
-   CSS (`.css`)
-   JavaScript (`.js`)
-   TypeScript (`.ts`)
-   YAML (`.yml`, `.yaml`)
-   Vue (`.vue`)

---

## Installation

Ensure you have **Go (version 1.18 or newer)** installed on your system.

You can install `commentpurger` with a single command:
```bash
go install github.com/scamgi/commentpurger@latest
```
This will compile the program and place the executable in your Go binary path (e.g., `$GOPATH/bin`), which should be part of your system's `PATH`.

---

## Usage

The basic command structure is to provide one or more paths (to files or directories) as arguments.

> **Warning:** This tool modifies your files directly. It's recommended to use it on files that are under version control.

#### Processing a Single File
```bash
commentpurger path/to/your/file.js
```

#### Processing an Entire Directory (and its subdirectories)
```bash
commentpurger ./my-project/
```

#### Processing Multiple Paths at Once
```bash
commentpurger ./src/components/ ./src/utils/main.js
```

### Example: Before and After

Imagine you have a JavaScript file `app.js`:

**Before `commentpurger`:**
```javascript
// app.js

// This is the main configuration object for our application.
const config = {
    apiKey: "ABC-123", /* IMPORTANT: Replace with a real key in production */
    featureFlags: {
        newUI: true, // Toggle this to switch to the new design
    }
};

/**
 * Initializes the application.
 * @param {object} initialConfig - The initial configuration.
 */
function initializeApp(initialConfig) {
    console.log("App is starting with config:", initialConfig);
    // TODO: Add more initialization logic here later.
}

initializeApp(config);
```

Run the command:
```bash
commentpurger app.js
```

**After `commentpurger`:**
```javascript
// app.js


const config = {
    apiKey: "ABC-123", 
    featureFlags: {
        newUI: true, 
    }
};


function initializeApp(initialConfig) {
    console.log("App is starting with config:", initialConfig);
    
}

initializeApp(config);
```

---

## Development

Want to contribute or build from the source yourself?

#### 1. Clone the Repository
```bash
git clone https://github.com/scamgi/commentpurger.git
cd commentpurger
```

#### 2. Build the Executable
```bash
go build
```
This will create the `commentpurger` binary in the current directory.

#### 3. Run Tests
To ensure everything is working as expected, run the test suite:
```bash
go test -v
```

---

## License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for the full text.