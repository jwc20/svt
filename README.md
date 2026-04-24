
<div id="top">

<!-- HEADER STYLE: CLASSIC -->
<div align="center">

<img src="https://github.com/user-attachments/assets/140838f6-b4d9-4744-ba4b-51dbd5285236" width="30%" style="position: relative; top: 0; right: 0;" alt="Project Logo"/>

# SVT

<em>Don't forget to pay the AWS bills!</em>

<!-- BADGES -->

<img src="https://img.shields.io/github/last-commit/jwc20/svt?style=default&logo=git&logoColor=white&color=0080ff" alt="last-commit">
<img src="https://img.shields.io/github/license/jwc20/svt?style=default&logo=opensourceinitiative&logoColor=white&color=0080ff" alt="license">

</div>

<br>

![demo](https://github.com/user-attachments/assets/43920e4b-c800-477c-80c0-004b97d7e66c)

<br>

---

## LIVE DEMO!

```bash
ssh siliconvalleytrail.xyz
```

---

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Quick Start](#quick-start)
  - [Installation](#installation)
  - [Running the Server](#running-the-server)
  - [Connecting via SSH](#connecting-via-ssh)
  - [Testing](#testing)
- [Project Structure](#project-structure)
- [Dependencies](#dependencies)


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

## Project Structure

```sh
└── svt/
    ├── README.md
    ├── cmd
    │   └── ssh              # SSH Wish Server
    ├── internal
    │   ├── engine           # Game Logic
    │   ├── hackernews       # Hacker News API Client
    │   ├── store            # SQLite Database
    └── └── ui               # Bubble Tea Terminal User Interface
```

---

## Dependencies

- [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) — Terminal user interface framework
- [Wish SSH](https://github.com/charmbracelet/wish) — Secure SSH server
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [HackerNews Algolia API](https://hn.algolia.com) — Bonus point

---

## See Also

- [otg](https://github.com/jwc20/otg)

---

<div align="right">

[![][back-to-top]](#top)

</div>

[back-to-top]: https://img.shields.io/badge/-BACK_TO_TOP-151515?style=flat-square
