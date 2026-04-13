# Game Deliverables and Requirements

## Core Requirements

1. **Game loop**
    - [x] A sequence of "days" or "turns."
    - [x] Each turn: choose an action with a **minimum of 3 options**, then resolve outcomes.
        - Actions: 
          1. Push forward 
          2. Fix a bug (play deathroll to decrease bug count)
          3. Marketing (spend cash to increase hype)

2. **Resources & state**
    - [x] Track at **least 3** meaningful stats. Examples:
        - Major resource:
          1. Cash
          2. Hype
          3. Tech Health
          
        - Other stats:
          1. Users Count 
          2. Bug Count
          3. Tech Debt
          4. Bonus Hype
   
    - [x] Winning happens when the team successfully reaches the destination.
      - Winning condition: reach 2040 miles 
      
    - [x] Define losing conditions
      - Losing major resource:
        1. Cash runs out (can't pay AWS bills)
        2. Hype runs out (lose from hype decay) 
        3. Tech Health runs out (lose from tech debt and bugs)

3. **Map**
    - [x] Represent progress across **at least 10** real, physical locations.
      - `1. San Jose → 2. Santa Clara → 3. Sunnyvale → 4. Mountain View → 5. Palo Alto → 6. Menlo Park → 7. Redwood City → 8. San Mateo → 9. Hillsborough → 10. San Bruno → 11. Daly City → 12. San Francisco`

4. **Events & choices**
    - [x] Implement a diverse set of events (some conditional on API data).
      - 20 unique events
    - [x] Provide consequential choices.
      - Player must choose between 3 actions, to carefully balance the resources to reach the destination. 
    - [x] An event should happen at each location after movement. Events should be at least semi-random, if not completely randomized.
      - At the end of each turn, a random event will occur based on the random.org number. Players have 30% chance of not getting an event.

5. **Public Web API integration (at least 1 required)**
    - [x] Use live or cached responses to **change gameplay** (not just display data).
      - [random.org](https://www.random.org/)
        - Used to randomize events and choices.
        - Used to determine winner of deathroll.
      - [HackerNews Algolia](https://hn.algolia.com)
        - Used to set bonus hype points based on player's name.  
    - [x] Provide a simple **fallback** (mock data) so the game runs without secrets or when offline.
      - If random.org fails, then it will use a mock number generator that I deployed on FastAPI cloud beta.
        - https://rng-api.fastapicloud.dev/random/100
        - If this fails, then it will use Go's built-in random number generator. (https://golang.org/pkg/math/rand/)
      - If HackerNews Algolia fails, then it will set bonus hype points to 0.

---


## Required Features

1. [x] **Testing:** at least a few unit tests covering critical logic (events, resource updates, win/lose).
    - Unit tests for the core game logic and random.org API client was written by me (`engine_test.go`, `rand_test.go`)
    - Extra unit and integration tests were written by Claude.
   
2. [ ] **Documentation:** a concise **README** (see Deliverables) and **Design Notes** explaining key choices and tradeoffs.
3. [ ] **Decency & safety:** handle API errors/timeouts gracefully, no hard-coded secrets, no collection of any personal user information.
   - No personal user information will be collected.
   - There are no secrets in the code.
   

## Deliverables

- [x] **Source code** in a public repo (any language/framework)
- [x] **Screen recording of or url to** the working application.
- **README.md** including:
    - [x] Quick start 
    - [ ] How to set API keys, how to run with mocks
      - No API keys are required to run the game.
    - [x] Brief architecture overview and dependency list
    - [x] How to run tests
    - [x] Example commands/inputs
    - [ ] How, if any, AI was utilized in the creation of the code and/or how the code utilizes AI, if at all

- **Design Notes** (can be a README section):
    - [x] Game loop & balance approach
    - [x] Why you chose your API(s) and how they affect gameplay
    - [ ] Data modeling (state, events, persistence)
    - [ ] Error handling (network failures, rate limits)
    - [ ] Tradeoffs and "if I had more time"
- [ ] **Tests** (unit or integration) for core mechanics