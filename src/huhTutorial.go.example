package main

import "github.com/charmbracelet/huh"
import "fmt"
import "errors"
import "log"

func main() {
    var (
        burger          string
        toppings        []string
        sauceLevel      int
        name            string
        instructions    string
        discount        bool
    )

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[string]().
                Title("Choose your burger").
                Options(
                huh.NewOption("Charmburger Classic", "classic"),
                huh.NewOption("Chickwhich", "chickwhich"),
                huh.NewOption("Fishburger", "fishburger"),
                huh.NewOption("Charmpossible Burger", "charmpossible"),
                ).
                Value(&burger),
            huh.NewMultiSelect[string]().
                Title("Toppings").
                Options(
                huh.NewOption("Lettuce", "lettuce").Selected(true),
                huh.NewOption("Tomatoes", "totmatoes").Selected(true),
                huh.NewOption("Jalapenos", "jalapenos"),
                huh.NewOption("Cheese", "cheese"),
                huh.NewOption("Vegan Cheese", "vegan cheese"),
                huh.NewOption("Nutella", "nutella"),
                ).
                Limit(4).
                Value(&toppings),
            huh.NewSelect[int]().
                Title("How much Charm Sauce do you want?").
                Options(
                huh.NewOption("None", 0),
                huh.NewOption("A little", 1),
                huh.NewOption("A Lot", 2),
                ).
                Value(&sauceLevel),
            ),

        huh.NewGroup(
            huh.NewInput().
                Title("What's your name?").
                Value(&name).
                Validate(func(str string) error {
                    if str == "Frank" {
                        return errors.New("Sorry, we don't serve customers named Frank.")
                    }
                    return nil
                }),
            huh.NewText().
                Title("Special Instructions").
                CharLimit(400).
                Value(&instructions),
            huh.NewConfirm().
                Title("Would you like 15% off?").
                Value(&discount),
            ),
        )

    err := form.Run()
    if err != nil {
        log.Fatal(err)
    }

    if !discount {
        fmt.Println("What? You didn't take the discount?!")
    }
}
