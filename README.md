# ktask

Ktask is a [kanban](https://en.wikipedia.org/wiki/Kanban_board) based todo-list
using a plain-text format to store data and use a TUI (terminal user interface)
to show the kanban board.

## Plain-text format
The plain-text format (as well as the parsing of it) is heavily inspired by the
[klog](https://github.com/jotaen/klog) project.

Essentially, it consists of a list of stages (for the kanban board). Each of
these stages then contain a list of task items in that stage. Each item has 3
properties: `createdAt` (when the entry was created), `modifiedAt` (when the
entry was last modified, aka moved to another stage) and `title`. These fields
must be present in that order. Similar to klog, the `title` might contain tags
starting with `#`, potentially storing a value. The first tag in the title will
be interpreted as the project to which this item belongs to.

### Example
```
todo
    2024-01-02 2024-02-01 buy milk #grocery
    2024-01-10 2024-02-11 Send out weekly newsletter #work=social

done
    2023-12-01 2024-01-01 celebrate new-year #friends
```
You may notice this example only consists of two stages instead the classical
three stages. This stresses, the file format does not impose any restrictions on
the order, number and name of the stages used.

## Project state
This project is still in a early stage. It can alreay be used (I'm doing so) but
it's not well tested and there might be bugs.

For open issues / planned features, see the
[issues page](https://github.com/atticus-sullivan/ktask/issues).

Although I want to avoid, there might be breaking changes (especially because of
the very early stage of the project). In order to get notified if there are any,
subscribe to [this issue](https://github.com/atticus-sullivan/ktask/issues/1).

## Acknowledgements
The basic idea for this project was greatly inspired by the charm tutorial
projects [`taskcli`](https://github.com/charmbracelet/taskcli) and
[`kancli`](https://github.com/charmbracelet/kancli/tree/main).

Thinking about data storage, the thought was usually there should not be that
much data (otherwise the board will be too full anyhow) so a plain-text file
should be enough (in contrast to a full blown database) while offering the
benefits of easy introspection and version control (if desired).

This lead to the process of thinking about the format. Since
[klog](https://github.com/jotaen/klog) has a similar approach (but for time
tracking) it only was logical to adapt its format and if so also make use of
its general parsing logic.
