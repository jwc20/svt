# Tradeoffs and If I Had More Time


## Tradeoffs

- Go vs. Python/JavaScript
  - Go was chosen because I wanted to create an SSH game using Bubble Tea and Wish.
  - I wanted to try something different rather than using the same tech stack as before. [bnc game](https://github.com/jwc20/bnc-game)

- CLI game vs. TUI game vs. Web (WebSocket) game

- WebSockets vs. ELM Architecture

- Spending four days learning the Oregon Trail BASIC code vs. diving into the SVT project immediately
  - I believe I made the right choice in studying the original game.
  - I wouldn't have known that the Oregon Trail code was one big while loop with a lot of if statements.
    - Understanding this helped me come up with my own ideas for the game.

## If I Had More Time

If given more time, I would focus on game design, balancing gameplay and mechanics, and studying newer versions of Oregon Trail.
I wanted to add more gameplay mechanics based on my own experience working at a four-person startup, but trying to implement them all in the game engine would have increased code complexity.
Many ideas written down in my journal were thrown away because they were not a good fit for the game engine or its balance.

**Examples:**

- [Adding difficulty levels using a weather API](journal/20260406.md)
- Resources and events that depend on date and weather
  - Idea: game is set before OpenAI releases ChatGPT (December 2022) [Major Event: ChatGPT releases](journal/20260408.md)
- Adjusting game formulas based on my experience working at a startup [Total expected revenue](journal/20260411.md)

---

In terms of additional features, it's hard to say. If I wanted a new feature (e.g. a new game action/event or a new mechanic), I could ask Claude Code to implement it and it would be done immediately.
The problem with hastily done, AI-written implementation is that it introduces code complexity and you lose understanding of your own codebase.

- Feature idea: view game history based on game state [SVT Notation](journal/20260406.md)
  - `[TurnNumber];[turn1:choice/eating]/[turn2:choice/eating]/.../....;[oxen/food/ammo/clothing/misc/skillLevel] - [event1/event2/.../...]`
  - View game history in the leaderboard