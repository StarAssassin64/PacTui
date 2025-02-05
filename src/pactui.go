package main

import (
    "bufio"
    "io"
    "log"
    "os"
    "os/exec"
    "runtime"
    "strings"
    "strconv"

    "github.com/charmbracelet/huh"
)

var (
    selectedMode    int
    // packages        string
    // search          string
    // results         []string
    selection       string
    numPackages     int
    pgCtr           int
    pgCount         int
    queryTitle      []string
)

func getNumPackages() (int, int) {
    command := "pacman -Q | wc -l"
    output, err := exec.Command("zsh", "-c", command).Output()
    if err != nil {
        log.Panicln(err.Error())
    }
    numPackagesInt := strings.Replace(string(output), "\n", "", -1)
    numPackages, err := strconv.Atoi(numPackagesInt)
    if err != nil {
        log.Panicln(err.Error())
    }
    numPages := getNumOfPages(numPackages)
    return numPackages, numPages
}

func clearScreen() {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("cmd", "/c", "cls")
    default:
        cmd = exec.Command("clear")
}
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func copyConf() []string{
    pmConfStat, err := os.Stat(os.ExpandEnv("/etc/pacman.conf"))
    if err != nil {
        log.Panicln("Fatal Error (-1) | Missing /etc/pacman.conf | Fix before using pactui")
    }
    if !pmConfStat.Mode().IsRegular() {
        log.Panicln("Fatal Error (-2) | /etc/pacman.conf is somehow not a regular file | Fix your installation of pacman before using pactui")
    }

    pmConf, err := os.Open(os.ExpandEnv("/etc/pacman.conf"))
    if err != nil {
        log.Panicln("Fatal Error (-3) | Unable to open /etc/pacman.conf")
    }

    defer pmConf.Close();

    ptConf, err := os.Create(os.ExpandEnv("$HOME/.config/pactui.conf"))
    if err != nil {
        log.Panicln("Fatal Error (-4) | Unable to create ~/.config/pactui.conf")
    }
    defer ptConf.Close()
    _, err = io.Copy(ptConf, pmConf)

    var returnValue []string
    scanning := bufio.NewScanner(pmConf)
    for scanning.Scan() {
        if strings.Contains(scanning.Text(), "#"){
            continue
        }
        returnValue = append(returnValue, scanning.Text())
    }

    return returnValue;

}

func confCheck() []string {
    var config []string
    _, err := os.Stat(os.ExpandEnv("$HOME/.config/pactui.conf"))
    if err != nil {
        config = copyConf()
    } else {
        ptConf, err := os.Open(os.ExpandEnv("$HOME/.config/pactui.conf"))
        if err != nil {
            log.Panicln("Fatal Error (-5) | Unable to open ~/.config/pactui.conf")
        }
        scanner := bufio.NewScanner(ptConf)
        for scanner.Scan() {
            if strings.Contains(scanner.Text(), "#"){
                continue
            }
            config = append(config, scanner.Text())
        }
    }
    return config
}

func getNumOfPages(numPackages int) int {
    return numPackages / 40
}

func getPage(page int) string {
    outputBytes, err := exec.Command("pacman", "-Q").Output()
    if err != nil {
        return err.Error()
    }
    output := string(outputBytes)

    outputLines := strings.Split(output, "\n")
    startingIndex := (page - 1) * 40
    endingIndex := startingIndex + 39
    log.Print(startingIndex, endingIndex)
    outputSlice := outputLines[startingIndex:endingIndex]
    return strings.Join(outputSlice, "\n")
}

func main() {
    // Default pacman.conf location: /etc/pacman.conf
    // Default pactui.conf location: ~/.config/pactui.conf
    config := confCheck()
    log.Print(config)

    mainPage := huh.NewGroup(
        huh.NewSelect[int]().
            Title("Welcome to PacTui").
            Options(
            huh.NewOption("Install a package", 0),
            huh.NewOption("Remove a package", 1),
            huh.NewOption("Search packages", 2),
            huh.NewOption("Query packages", 3),
            huh.NewOption("Configurate PacTui", 4),
            huh.NewOption("Exit PacTui", 5),
            ).
            Value(&selectedMode),
        )

    clearScreen()
    pgCtr = 1
    application :=huh.NewForm(mainPage)
    application.Run()
    switch selectedMode {
    case 3:
       queryPageRun()
    case 5:
    default:
        os.Exit(0)
}

}

func queryPageRun() {
    numPackages, pgCount = getNumPackages()
    queryTitle = []string{ "Query (" , strconv.Itoa(numPackages) , "Packages | Page", strconv.Itoa(pgCtr), "/", strconv.Itoa(pgCount), "):" }
    clearScreen()

    queryPage := huh.NewGroup(
        huh.NewNote().
            Title(strings.Join(queryTitle , " ")).
            Height(40).
            Description(getPage(pgCtr)),
        huh.NewSelect[string]().
            Options(
            huh.NewOption("Prev Page", "pg-"),
            huh.NewOption("Home", "home"),
            huh.NewOption("Next Page", "pg+"),
            ).
            Value(&selection),
        )

    application := huh.NewForm(queryPage)
    application.Run()
    if selection == "home" {
        main()
    }
    if selection == "pg+" {
        pgCtr ++
        if pgCtr > pgCount {
            pgCtr--
            queryPageRun()
        }
        queryPageRun()
    } else if selection == "pg-" {
        pgCtr --
        if pgCtr <= 0 {
            pgCtr++
            queryPageRun()
        }
        queryPageRun()
    }
}
