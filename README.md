# ebiten-holiday-2024
### Rhythm Game
A rhythm game that utilizes "areas" of the keyboard as single inputs. These areas are composed of sets of adjacent keys. 
The game is played by pressing the keys in the area that corresponds to the incoming notes in time with the music.

The goal of the game is achieve the highest score possible by hitting the notes accurately and consistently. Players are graded on their performance at the end of each song and receive a rank based on their score.

Players are able to select from a list of songs to play, each with two difficulties. The game features a tutorial to help new players learn the mechanics of the game.

The game also features a settings menu where players can adjust various settings such as keybinds, audio settings, and video settings.

The game also has a chart editor that allows players to create their own charts for songs. Players can import audio, edit metadata, create and edit charts, place notes, playtest charts, and export song packages that can be shared.

## Features
### To Be Determined
- Name
- Theme Integration
- The Game Inputs
    - Five quadrants:
        - ` -> 6 , Tab -> T
        - Capslock -> G, Left Shift -> B 
        - 7 -> Backspace, Y -> \
        - H -> Enter, N -> Right Shift
        - Space Bar
    - Two more quadrants: (?)
        - Arrow Keys
        - Nav Cluster (the keys above the arrow keys)
    - Six individual keys (?)
        - Use center row of keyboard
        - A S D, L ; '
- The Note + Track representation for base renderer
    - Circular?
    - Linear?

### Need
- Songs (aiming for at least 6-8)
    - Packaged
    - Metadata
    - Audio
    - Art
    - Charts
        - Two difficulties (due to time constraints)
            - easy
            - hard
- Settings
    - Keybinds
        - Area Mapping
        - ...
    - Audio
        - Latency Adjustment
        - Volume
        - ...
    - Video
        - Renderer
        - Resolution
        - Fullscreen
        - ...
    - Game
        - Note Speed
        - ...
- Game
    - Tutorial
    - Note spawning
    - Note hit detection
    - Note travel
    - Input handling
    - Scoring
    - Feedback
        - Hit audio
        - Hit animations
        - Miss animations
        - Score display
        - Combo display
    - Health/Fail State (? TBD)
    - Results
        - Score
        - Combo
        - Accuracy
        - Grade
        - Rank
        - Leaderboard (? TBD)
- Systems
    - Audio
        - Playback
        - Latency Adjustment
        - Syncing
    - Chart Editor
        - Import audio
        - Edit metadata
        - Create and edit charts
        - Place notes
        - Playtest charts
        - Export song package
    - Input
        - Intercept all keyboard input
        - Map keys to areas
        - Handle multiple keys pressed at once
    - Rendering
        - Interchangeable (2d / 3d / etc)
        - Divorced from game logic

### Nice (no particular order (server's probably nicest))
- Server
    - Hosts songs/leaderboard/profiles
    - In-game direct song downloads
    - Persistent User Profile
    - Online Leaderboards
- More songs
- UI Skins
- Song Filter/Order

## UI Brainstorming
- Intro
    - Natural transition into Main Menu
- Main Menu
    - "Play" (Song Selection)
        - Perform latency check if not set (? maybe can build it into tutorial)
        - Prompts to play tutorial on first launch
    - Settings
    - Credits
    - Exit
- Song Selection
    - Song List
        - Wheel?
        - Flat List?
    - Settings
    - Exit
- Song (once highlighted in Song Selection)
    - Song Audio Preview
    - Details Render
        - Album Art
        - Title
        - Artist
        - Difficulty
        - Length
        - Notes
        - BPM
        - Charted By
        - High Score
        - Note Heatmap (?)
    - Difficulty Selection
    - Play