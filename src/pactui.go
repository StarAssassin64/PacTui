package main

import (
    "bufio"
    "io"
    "log"
    "os"
    "os/exec"
    "runtime"
    "strconv"
    "strings"

    "github.com/charmbracelet/huh"
)

var (
    selectedMode     int
    packageInput     string
    selectedPackages []string
    selection          string
    numPackages        int
    pgCtr              int
    pgCount            int
    queryTitle         []string
    visability         bool
    errorTitle         string
    errorDesc          string
    packageCheckString string
    packageInstallLog  []byte
	packageRemovalLog  []byte
    confTitle          string
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

func copyConf() []string {
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

    defer pmConf.Close()

    ptConf, err := os.Create(os.ExpandEnv("$HOME/.config/pactui.conf"))
    if err != nil {
        log.Panicln("Fatal Error (-4) | Unable to create ~/.config/pactui.conf")
    }
    defer ptConf.Close()
    _, err = io.Copy(ptConf, pmConf)

    var returnValue []string
    scanning := bufio.NewScanner(pmConf)
    for scanning.Scan() {
        if strings.Contains(scanning.Text(), "#") {
            continue
        }
        returnValue = append(returnValue, scanning.Text())
    }

    return returnValue

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
            if strings.Contains(scanner.Text(), "#") {
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
            huh.NewOption("Remove a package [Not Implemented]", 1),
            huh.NewOption("Search packages [Not Implemented]", 2),
            huh.NewOption("Query packages", 3),
            huh.NewOption("Configure PacTui [Not Implemented]", 4),
            huh.NewOption("Exit PacTui", 5),
            ).
            Value(&selectedMode),
    )

    clearScreen()
    pgCtr = 1
    application := huh.NewForm(mainPage)
    application.Run()
    switch selectedMode {
	case 1:
		writeRemovalPageRun()
    case 3:
        queryPageRun()
    case 0:
        writeInstallPageRun()
    case 5:
    default:
        os.Exit(0)
    }

}

func queryPageRun() {
    numPackages, pgCount = getNumPackages()
    queryTitle = []string{"Query (", strconv.Itoa(numPackages), "Packages | Page", strconv.Itoa(pgCtr), "/", strconv.Itoa(pgCount), "):"}
    clearScreen()

    queryPage := huh.NewGroup(
        huh.NewNote().
            Title(strings.Join(queryTitle, " ")).
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
        pgCtr++
        if pgCtr > pgCount {
            pgCtr--
            queryPageRun()
        }
        queryPageRun()
    } else if selection == "pg-" {
        pgCtr--
        if pgCtr <= 0 {
            pgCtr++
            queryPageRun()
        }
        queryPageRun()
    }
}

func validatePackages(packageString string) bool {
    packageValArr := strings.Split(packageString, " ")
    for i := 0; i < len(packageValArr); i++ {
        packageCheckString = "^" + packageValArr[i] + "$"
        output, err := exec.Command("pacman", "-Ssq", packageCheckString).Output()
        if err != nil {
            return false
        }
        if string(output) == "" || output == nil {
            return false
        }
    }
    return true
}

func validateRemoval(packageString string) bool {
	_, err := exec.Command("pacman", "-Qq", packageString).Output()
	if err != nil {
		return false
	} else {
		return true
	}
}

func writeInstallPageRun() {
    writeInstallPage := huh.NewGroup(
        huh.NewNote().
            TitleFunc(func() string {
                if visability {
                    errorTitle = "ERROR:"
                } else {
                    errorTitle = ""
                }
                return errorTitle
            }, &errorTitle).
            DescriptionFunc(func() string {
                if visability {
                    errorDesc = "One or more of the packages does not exist, please check spelling"
                } else {
                    errorDesc = ""
                }
                return errorDesc
            }, &errorDesc),
        huh.NewInput().
            Title("Install Packages").
            Value(&packageInput),
        huh.NewSelect[string]().
            Options(huh.NewOption("Home", "home"),
            huh.NewOption("Install", "install"),
            ).
            Value(&selection),
    )

    application := huh.NewForm(writeInstallPage)
    application.Run()
    if !validatePackages(packageInput) {
        visability = true
        writeInstallPageRun()
    }
    if selection == "home" {
        main()
    } else if selection == "install" {
        confTitle = packageInput
        installPageRun()
    }

}

func installPageRun() {

    confirmationPage := huh.NewGroup(
        huh.NewSelect[string]().
            Title("Are you sure you want to install "+confTitle).
            Options(
            huh.NewOption("Yes", "y"),
            huh.NewOption("No", "n"),
            ).Value(&selection),
    )

    successfulinstallPage := huh.NewGroup(
        huh.NewSelect[string]().
            Title("Successfully Installed Packages").
            Options(
            huh.NewOption("Home", "home"),
            huh.NewOption("Back", "back"),
            ).
            Value(&selection),
    )
    failedInstallPage := huh.NewGroup(
        huh.NewNote().
            Title("Failed Installing Packages").
            Description(string(packageInstallLog)).
            Height(40),
        huh.NewSelect[string]().
            Options(
            huh.NewOption("Home", "home"),
            huh.NewOption("Back", "back"),
        ),
    )

    application := huh.NewForm(confirmationPage)
    application.Run()
    if selection == "y" {
        returnValue := installPackages(strings.Split(confTitle, " "))
        if returnValue {
            application = huh.NewForm(successfulinstallPage)
            application.Run()
            if selection == "home" {
                main()
            }
            if selection == "back" {
                writeInstallPageRun()
            }
        } else {
            application = huh.NewForm(failedInstallPage)
            application.Run()
            if selection == "home" {
                main()
            }
            if selection == "back" {
                writeInstallPageRun()
            }
        }
    }

}

func installPackages(packages []string) bool {
    var err error
    packageInstallLog, err = exec.Command("pacman", "-Sy", "--noconfirm", strings.Join(packages, " ")).Output()
    println(string(packageInstallLog))
    if err != nil {
        return false
    }
    packageInstallLog, err = exec.Command("tail", "-40", string(packageInstallLog)).Output()
    println(string(packageInstallLog))
    return true
}

func writeRemovalPageRun() {
    writeRemovalPage := huh.NewGroup(
        huh.NewNote().
            TitleFunc(func() string {
                if visability {
                    errorTitle = "ERROR:"
                } else {
                    errorTitle = ""
                }
                return errorTitle
            }, &errorTitle).
            DescriptionFunc(func() string {
                if visability {
                    errorDesc = "One or more of the packages does not exist, please check spelling"
                } else {
                    errorDesc = ""
                }
                return errorDesc
            }, &errorDesc),
        huh.NewInput().
            Title("Remove Package").
            Value(&packageInput),
        huh.NewSelect[string]().
            Options(huh.NewOption("Home", "home"),
            huh.NewOption("Remove", "remove"),
            ).Value(&selection),
    )

    application := huh.NewForm(writeRemovalPage)
    application.Run()
    if !validateRemoval(packageInput) {
        visability = true
        writeRemovalPageRun()
    }
    if selection == "home" {
        main()
    } else if selection == "remove" {
        confTitle = packageInput
        removalPageRun()
    }
}

func removalPageRun() {
    confirmationPage := huh.NewGroup(
        huh.NewSelect[string]().
            Title("Are you sure you want to remove "+confTitle).
            Options(
            huh.NewOption("Yes", "y"),
            huh.NewOption("No", "n"),
            ).Value(&selection),
    )

    successfulRemovalPage := huh.NewGroup(
        huh.NewSelect[string]().
            Title("Successfully Removed Package").
            Options(
            huh.NewOption("Home", "home"),
            huh.NewOption("Back", "back"),
            ).
            Value(&selection),
    )
    failedRemovalPage := huh.NewGroup(
        huh.NewNote().
            Title("Failed Removing Package").
            Description(string(packageRemovalLog)).
            Height(40),
        huh.NewSelect[string]().
            Options(
            huh.NewOption("Home", "home"),
            huh.NewOption("Back", "back"),
        ),
    )

    application := huh.NewForm(confirmationPage)
    application.Run()
    if selection == "y" {
        returnValue := removePackage(confTitle)
        if returnValue {
            application = huh.NewForm(successfulRemovalPage)
            application.Run()
            if selection == "home" {
                main()
            }
            if selection == "back" {
                writeRemovalPageRun()
            }
        } else {
            application = huh.NewForm(failedRemovalPage)
            application.Run()
            if selection == "home" {
                main()
            }
            if selection == "back" {
                writeRemovalPageRun()
            }
        }
    }

}

func removePackage(packageString string) bool {
	log, err := exec.Command("pacman", "-Ruv", packageString, "--noconfirm").Output()
	packageRemovalLog = log
	if err != nil {
		return false
	} else {
		return true
	}
}
