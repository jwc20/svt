<div id="top">

<!-- HEADER STYLE: CLASSIC -->
<div align="center">

<img src="notes/img/Sprite-silicon-stone.png" width="30%" style="position: relative; top: 0; right: 0;" alt="Project Logo"/>

# Silicon Trail

<em>Don't forget to pay the AWS bills!</em>

<!-- BADGES -->

<img src="https://img.shields.io/github/last-commit/jwc20/svt?style=default&logo=git&logoColor=white&color=0080ff" alt="last-commit">


</div>

<br>

![demo](https://github.com/user-attachments/assets/43920e4b-c800-477c-80c0-004b97d7e66c)

<br>

---

## LIVE DEMO!

You can run ssh into the deployed server by running the command:

```bash
ssh siliconvalleytrail.xyz

# if you don't have public key then run the command below:
# ssh-keygen -t ed25519 -C "your_email@email.com"
```


## Deliverables

- See [deliverables-and-req.md](notes/deliverables-and-req.md)

## Game Design

- See [game-design-doc.md](notes/game-design-doc.md)

## AI Uses

- See [ai-uses.md](notes/ai-uses.md)

## Tradeoffs & If I had more time

- See [tradeoffs-and-if-i-had-more-time.md](notes/tradeoffs-and-if-i-had-more-time.md)

## API Choices

- [random.org](https://www.random.org/)
- [HackerNews Algolia](https://hn.algolia.com)

---

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Quick Start](#quick-start)
  - [Installation](#installation)
  - [Running the Server](#running-the-server)
  - [Connecting via SSH](#connecting-via-ssh)
  - [Testing](#testing)
- [Dependencies](#dependencies)
- [Project Structure](#project-structure)


---

## Quick Start

### Installation

Clone the repository and set up the project:

1. **Clone the repository:**

```sh
git clone https://github.com/jwc20/svt
```

2. **Navigate to the project directory:**

```sh
cd svt
```

3. **Build the application:**

Using [Go modules](https://go.dev/blog/using-go-modules):

```sh
go build -o svt cmd/ssh
```

### Running the Server

**Option 1: Run the built binary**

```sh
./svt
```

**Option 2: Run directly without building**

```sh
go run ./cmd/ssh
```

### Connecting via SSH

After starting the server, connect to it using SSH:

```sh
ssh localhost -p 23234

# or with a username
ssh player_name@localhost -p 23234
```

https://github.com/user-attachments/assets/ec4471b6-fd84-4718-ad35-6a4e1ba8b411

### Testing

Run the test suite with:

```sh
go test ./...
```

---

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — Terminal user interface framework
- [Wish SSH](https://github.com/charmbracelet/wish) — Secure SSH server
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — UI components (tables, viewports)
- [joho/godotenv](https://github.com/joho/godotenv) — Environment variable loading
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) — SQLite driver
- [stretchr/testify](https://github.com/stretchr/testify) — Testing assertions

---


## Project Structure

```sh
└── svt/
    ├── README.md
    ├── cmd
    │   └── ssh              # SSH Wish Server
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
