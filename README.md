# GoHexaGen: Clean Architecture Scaffolder for Go 🐹🛡️

"Implicit boundaries. Unlimited performance."

[![Architecture: Hexagonal](https://img.shields.io/badge/Architecture-Hexagonal-blueviolet)](#)
[![Go: 1.21+](https://img.shields.io/badge/Go-1.21%2B-blue)](#)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](#)

---

## 🚀 Why GoHexaGen?

Go applications thrive on speed and minimalism. However, managing clean architecture interfaces can sometimes feel like a manual task. **GoHexaGen** provides the essential scaffolding support to decouple components dynamically.

### Key Features
- **Primary Adapters**: Support for **Fiber**, **Gin**, and **Echo**.
- **Implicit Ports**: True decoupling through standard Go interfaces.
- **Graceful Shutdowns**: OS Signal interception for safe teardowns.
- **Enterprise Logging**: Standard `log/slog` structured reporting.

## 📖 Quick Start

### 1. Install Globally
To install the CLI directly on your machine, run:
```bash
go install github.com/rbaezc/gohexagen@latest
```

### 2. Initialize a New Project
Run the generator interactively:
```bash
gohexagen init my_hexagonal_service
```

### 3. Generate a Domain Slice
Generate controllers, traits, models, and DB queries in one command:
```bash
cd my_hexagonal_service
gohexagen gen-resource User name:string email:string
```

### 4. Run Tests
```bash
go test ./...
```

Developed by **rbaezc**. Licensed to **Vortex Solutions**.
