## School Project ##
This project was created for a class called Computer Science Projects

This project is graded

## Project Description ##
This is going to be a Terminal User Interface (TUI) front end for [pacman](https://gitlab.archlinux.org/pacman/pacman) (Arch Linux's package manager)

## Planned options ##
- Install [Implemented]
    - Install packages listed in the text box
    - Provide a log after entering for the success or failure of the installation of the packages
- Remove
    - Remove packages listed in the text box
    - Provide a log after entering for the success or failure of the removal of the packages
- Query [Implemented]
    - List all installed packages in a user readable format
- Search
    - Search for a package in a user readable format
    - Provide a simple way to install a selected package 
- Config
    - Configure the pacman.conf for PacTui
    - Provide an easy way to find the best mirrors to use and select which mirrors to use after finding

## Other Details ##
- Written in GO using [Huh?](https://github.com/charmbracelet/huh) as a library
