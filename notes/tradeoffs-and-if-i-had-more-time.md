# Tradeoffs and if I had more time


## Tradeoffs

- Go vs. Python/Javascript
    - Go was chosen because I wanted to create a SSH game using Bubble Tea and Wish.
    - I wanted to do something different rather than using the same tech stack as before. [bnc game](https://github.com/jwc20/bnc-game)

- CLI game vs. TUI game vs. Web (Websocket) game 

- Websockets vs. ELM Architecture

- Spending four days learning Oregon Trail BASIC code vs. Diving into SVT project immediately
  - I believe I made the right choice in studying the original game.
  - I wouldn't have know the Oregon Trail code was one big while loop with a lot of if statements.
    - Understanding this helped me coming up with my own ideas for the game.

## If I had more time

If given more time, I would spend more time on game design, balancing gameplay and gameplay mechanics, and learning more about the game of newer versions of Oregon Trail.
I wanted to add more gameplay mechanics, based on my own experience working at a four person startup. But the trying to implement them all in the game engine would have increased the code complexity.
Many ideas that was written down in the journal was thrown away because it was not a good fit for the game engine and balance.

**Examples:**

- [Adding Difficulty Levels, using weather API](notes/journal/20260406.md)
- Resources and events depends on date and weather
  - Idea: game sets before OpenAI releases ChatGPT (December 2022) [Major Event: ChatGPT releases](notes/journal/20260408.md)
- Adjust game formulas based on my experience working at a startup. [Total expected revenue](notes/journal/20260411.md)

---

In terms of additional features, it's hard to say. Because if I wanted a new feature (e.g. a new game action/event, or a new game mechanic), I could ask Claude Code to implement it and it would be done immediately.
The problem with hastily-done/AI-written implementation is code quality and complexity and lose understanding of your own code base.

- Feature idea: able to see game history based on game state [SVT Notation](notes/journal/20260406.md)
  - `[TurnNumber];[turn1:choice/eating]/[turn2:choice/eating]/.../....;[oxen/food/ammo/clothing/misc/skillLevel] - [event1/event2/.../...]`
  - See game history in the leaderboard