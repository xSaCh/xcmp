# XCmp

A minimal TUI client for MPRIS that displays song title, artist, progress, and album art.

![xcmp](https://github.com/user-attachments/assets/0dd4b5b9-737d-42c9-8e42-09461c66602b)

- [BubbleTea](https://github.com/charmbracelet/bubbletea) for TUI.
- `MPRIS` and `DBUS` for song streaming information.
- [Kitty Graphic Protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/) for drawing images.

**NOTE:** Works only in Linux with Kity or Terminal emulators supporting Kitty Graphic Protocol such as Konsole.

## Installation

1. **Clone the repository:**
   ```sh
   $ git clone https://github.com/xSaCh/xcmp
   $ cd xcmp
   ```
2. **Build the application:**
   ```sh
   $ go build -o xcmp
   ```

## Usage

List available MPRIS players:
```sh
$ ./xcmp list-players
```

Specify an MPRIS player:
```sh
$ ./xcmp -player "player_name"
```

Get a list of available players using the `list-players` command.
