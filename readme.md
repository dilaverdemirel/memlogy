# Memlogy

Memlogy is a daily notebook that you can easily add note, terminal app.

![list](/docs/list-1.png)

I tried something for learning golang. It contatains;
- github.com/eiannone/keyboard : for handle keyboard inputs
- github.com/mattn/go-sqlite3 : for store notes
- github.com/olekukonko/tablewriter : for list notes
- github.com/spf13/cobra : for CLI
- github.com/spf13/viper : for configuration

## Usage 

### Add new note :

```sh
memlogy add Write something about you
memlogy add --description "Write something about you"
memlogy add -d "Write something about you"

#Adds entry to yesterday
memlogy add --go-to-day=-1 --description="Write something about you" #For yesterday with --go-to-day flag
memlogy add -g=-1 -d="Write something about you" #For yesterday with -g flag
```

### List notes :

```sh
memlogy list
```

You can use keyboard for:
- **Left Arrow** go to previous days
- **Right Arrow** go to next days
- **Esc** to exit list
- **Delete** to delete an item

![delete](/docs/list-2.png)