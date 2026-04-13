<div id="top">

<!-- HEADER STYLE: CLASSIC -->
<div align="center">

<img src="notes/journal/img/Sprite-silicon-stone.png" width="30%" style="position: relative; top: 0; right: 0;" alt="Project Logo"/>

# Silicon Trail

<em>Don't forget to pay the AWS bills!</em>

<!-- BADGES -->

<img src="https://img.shields.io/github/last-commit/jwc20/svt?style=default&logo=git&logoColor=white&color=0080ff" alt="last-commit">


</div>

<br>

![demo](https://github.com/user-attachments/assets/43920e4b-c800-477c-80c0-004b97d7e66c)

<br>

---

## Deliverables

- See [deliverables-and-req.md](notes/deliverables-and-req.md)

## Game Design

- See [game-design-doc.md](notes/game-design-doc.md)

## AI Uses

- See [ai-uses.md](notes/ai-uses.md)

## Tradeoffs & If I had more time

- See [tradeoffs-and-if-i-had-more-time.md](notes/tradeoffs-and-if-i-had-more-time.md)

### API Choices

- [random.org](https://www.random.org/)
- [HackerNews Algolia](https://hn.algolia.com)

<br>

---


## Table of Contents

- [Table of Contents](#table-of-contents)
- [Quick Start](#quick-start)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Testing](#testing)
- [Dependencies](#dependencies)
- [Project Structure](#project-structure)


---

## Quick Start

### Installation

Build svt from the source and install dependencies:

1. **Clone the repository:**

```sh
git clone https://github.com/jwc20/svt
```

2. **Navigate to the project directory:**

```sh
cd svt
```

3. **Install:**

**Using [go modules](https://golang.org/):**

```sh
go build -o svt cmd/ssh # build the binary
./svt                   # run
```

### Usage

Run the project with:

```sh
# If the binary is built
go build -o svt cmd/ssh 
./svt

# If the binary is not built
go run ./cmd/ssh
```

After running the Wish server, you can connect to it using SSH:

```sh
ssh localhost -p 23234

# or
ssh player_name@localhost -p 23234
```


https://github.com/user-attachments/assets/ec4471b6-fd84-4718-ad35-6a4e1ba8b411


### Testing

```sh
go test ./...
```

---

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) 
- [Wish SSH](https://github.com/charmbracelet/wish)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [Bubbles](https://github.com/charmbracelet/bubbles) for tables and viewports
- [joho/godotenv](https://github.com/joho/godotenv) for loading environment variables
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) SQLite driver for Go
- [stretchr/testify](https://github.com/stretchr/testify) for writing asserts

---


## Project Structure

This application uses the [Bubble Tea v2 framework](https://github.com/charmbracelet/bubbletea) for terminal user interfaces and the [lipgloss](https://github.com/charmbracelet/lipgloss) library for styling.

It uses [Wish](https://github.com/charmbracelet/wish) to provide a secure SSH connection to the server.

The core game logic and API client is in the `internal/engine`, `internal/rand`, and `internal/hackernews` directories.

```sh
└── svt/
    ├── README.md
    ├── cmd
    │   └── ssh              # SSH Wish Server
    ├── go.mod
    ├── go.sum
    ├── internal
    │   ├── engine           # Game Logic
    │   ├── rand             # Random Number Generation API Client
    │   ├── hackernews       # Hacker News API Client
    │   ├── store            # SQLite Database
    │   └── ui               # Bubble Tea Terminal User Interface
    └── notes
```



<div align="right">

[![][back-to-top]](#top)

</div>

[back-to-top]: https://img.shields.io/badge/-BACK_TO_TOP-151515?style=flat-square

