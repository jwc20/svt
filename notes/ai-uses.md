# AI uses


## Day 0-4

- To figure out the game flow and gameplay logic for Oregon Trail, I built the game for the CLI based on the original BASIC game and API client for random.org.
  - [My oregon trail cli repo](https://github.com/jwc20/oregon-trail-go)
  - [The original BASIC game and Python Version](https://github.com/philjonas/oregon-trail-1978-python)

- Very little AI was used during this stage, when used, it was to translate some of the confusing BASIC code.
- AI was not used to brainstorming gameplay/feature ideas.
- Github Desktop's built-in AI was used to create commit messages throughout the project.
  - ![ai-use1.jpg](journal/img/ai-use1.jpg)

## Day 5

- Claude Opus (chat) and Claude Code was used: 
  - to convert CLI game to Bubble Tea/Wish TUI game (used mostly for game TUI styling and layout),
  - [to refactor and reorganize the code structure](notes/journal/20260410.md)

- Game prompts were written using Claude Opus.

![ai-use0.jpg](journal/img/ai-use0.jpg)

## Day 6-7

- Claude Opus (chat) was used to organize journal/ideas files and pen and paper notes.

- After finishing building the core game logic and TUI implementation, Claude Code was used to create additional features.

  - [Implement SQLite persistence and enhance game state management (#5)](https://github.com/jwc20/svt/commit/d976eaec2ec28119e1deda0d65f7de8b37f9a09b)
  - [Implement leaderboard UI (#6)](https://github.com/jwc20/svt/commit/f542284a9866b595c25f70b17ef9c45be4655573)
  - HackerNews Algolia API Client
  - Third action (Marketing Push)

- Tests that is not in `engine_test.go` and `rand_test.go` were written using Claude Code.
