# Level Editor

The level editor uses ASCII characters to define levels.

## Structure

The level templates contain three sections

* CONFIG
* HEADER
* LEVEL

### Config

The config section is used to specify what modifiers are present on different types of symbols. For
example, `%: Color:Green` makes the player green.

### Configuration Language
The configuration values are split using the following syntax: `characters: modifiers` Multiple characters and modifiers may be added.
For example: `cud: Color:Yellow` would color the `c`, `u`, and `d` buttons yellow.

#### Color
The `Color` modifier allows you to specify the color as a parameter using the following syntax: `Color:COLOR`.
For example, `Color:Yellow`.

### Header

The header section is used to render a header in the terminal when you're on that level. For example, level_0 introduces
the game

```
HEADER
################################
##        COLORMANCER         ##
################################
```

### Level

The level section is used to define the actual level.

## Definitions

| Character | Description  | Allowed Modifiers   |
|-----------|--------------|---------------------|
| @         | The Player   | Color               |
| #         | Walls        | Color               |
| A-Z       | Door         | Color, Open, Closed |
| a-z       | Button       | Color               |
| 0-2       | Color Pickup | N/A                 |
| %         | Exit         | Color               |

### Special Interactions
* `A-Z` doors are triggered by the corresponding `a-z` button. The `Open/Closed` modifiers specify the initial state, and the buttons toggle to the opposite state.
* The `0-2` color pickups are `Cyan`, `Red`, `Yellow`. They do not accept modifiers
* The color of `Exits` & `Buttons` determines what color must interact with it for it to be activates

### Additional Remarks
While possible, please do not change colors of walls.

