# Twitch Plays Into The Breach

* Run the game at the top left of the screen.
* Set up your stream
* Start the app and hit OK when the game is at the forefront
* You're good to go

## Configuration

Set the following environment variables:

- `TWITCH_USERNAME` - lowercase Twitch username of the Bot
- `TWITCH_TOKEN` - OAuth token of the Bot, starting with `oauth:`
- `TWITCH_CHANNEL` - Stream channel name

## Chat commands

Mouse:
- `click` - Click the mouse
- `mouse [x] [y]` - Move mouse to given x, y pixel coordinates
- `click [x] [y]` - Click at a given x y point (similar to doing `mouse x y` then `click`)
- `mouse [xy]` - Mouse at map coordinate
- `click [xy]` - Click at map coordinate

Mechs and units:
- `select [mech|deployed|mission] [1-3]` - Select mech/deployed/mission unit 1 to 3
- [ ] `move [mech|deployed|mission] [1-3] [xy]` - Move unit to a given map coordinate
- `weapon [1-2]` - Select weapon
- `deselect|disarm` - Deselect / disarm weapon
- `attack [mech|deployed|mission] [1-3] [weapon 1-2] [xy]` - Attack with a mech using a weapon at a given map coordinate
- `repair [mech 1-3] [xy]` - Repair the mech at this map coordinate
- `next` - Select next unit. Also use `next 2` to select the second next unit and so on

Moves and turns:
- `undo` - Undo move
- `reset` - Reset turn
- `endturn` - End turn

UI:
- `info [on|off]` - Toggle the info overlay
- `order [on|off]` - Toggle the attack order overlay