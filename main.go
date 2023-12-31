package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

func main() {
	tpl := template.Must(template.ParseFiles("templates/basic/test.gohtml"))
	_ = tpl.Execute(os.Stdout, "go")

	fmt.Println("------------------")
	// when the data attribute is structured its values
	// can be accessed by calling the dot operator
	// with its exported fields
	tpl = template.Must(template.ParseFiles("templates/basic/testStruct.gohtml"))
	_ = tpl.Execute(os.Stdout, struct {
		Name string
		Age  int
		Mail string
	}{
		Name: "Walter",
		Age:  32,
		Mail: "walter@mail.com",
	})

	fmt.Println("------------------")
	// Templates parts can be executed conditionally by the use of
	// {{ if . }} statements with its pipelines.
	// Alongside the if operator its possible to use {{ else }}
	// and {{ else if . }}

	tpl2 := template.Must(template.ParseFiles("templates/basic/testConditional.gohtml"))
	_ = tpl2.Execute(os.Stdout, true)
	_ = tpl2.Execute(os.Stdout, false)

	fmt.Println("------------------")
	// Range makes able to iterate trough iterators(!)
	// and if needed access the data stored inside of them
	// with the dot

	tpl3 := template.Must(template.ParseFiles("templates/basic/testRange.gohtml"))
	_ = tpl3.Execute(os.Stdout, []string{"test1", "test2", "test3"})

	fmt.Println("------------------")
	// When using Range with an iterable where the value is a struct is possible
	// to easily access its properties with the dot operator

	// Alongside the range operator it is possible to use
	// {{break}} and {{continue}} for better control of
	// execution flow
	tpl = template.Must(template.ParseFiles("templates/basic/testRangeStruct.gohtml"))

	type dog struct {
		Name  string
		Breed string
		Age   int
	}

	dogs := []dog{{"luna", "Yorkie", 2}, {"nina", "Border Collie", 14}}
	_ = tpl.Execute(os.Stdout, dogs)

	fmt.Println("------------------")
	// Ranges structs and conditional can all be
	// used together to elaborate how the output of the
	// evaluated template will be depending on its input data

	tpl = template.Must(template.ParseFiles("templates/basic/testRangeConditionalChain.gohtml"))
	_ = tpl.Execute(os.Stdout, dogs)

	fmt.Println("------------------")
	// It's possible to extend the capabilities of the template engine by injecting functions
	// to be used with pipelines

	//Map of functions to be passed to the template,
	//so they can be used when evaluating the template white `Execute`
	funcMap := map[string]interface{}{
		"toUpper": strings.ToUpper,
		"customFunction": func(input string) string {
			return input + "!"
		},
		// Note that when applying the function with mismatched types
		// execution does not break but the output will be omitted
		"double": func(input float64) float64 {
			return input * 2
		},
	}

	tplInjected := template.Must(template.New("testFunctionInjection.gohtml").
		Funcs(funcMap).
		ParseFiles("templates/basic/testFunctionInjection.gohtml"))
	// Note that to use `Funcs` is needed to already have allocated the template wit `New`
	// and for New` to work with `ParseFiles` the name *must* be equal to the file name being parsed

	_ = tplInjected.Execute(os.Stdout, "text")

	fmt.Println("------------------")
	// Creating a basic template composition by defining 3 templates(header, body and footer)
	// and using them in a separated `composition.gohtml` template without any internal data

	prefix := "templates/composition/basic/"
	files := []string{
		prefix + "header.gohtml",
		prefix + "body.gohtml",
		prefix + "footer.gohtml",
		prefix + "composition.gohtml",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	_ = tpl.ExecuteTemplate(os.Stdout, "composition.gohtml", nil)

	fmt.Println("------------------")
	// Creating another template composition but using data injection with defined type
	//structs for each composited template.

	prefix = "templates/composition/advanced/"
	files = []string{
		prefix + "header.gohtml",
		prefix + "body.gohtml",
		prefix + "footer.gohtml",
		prefix + "composition.gohtml",
	}
	tpl, err = template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}

	type Header struct {
		Title string
	}

	type Body struct {
		Title string
		Text  string
	}

	type Footer struct {
		Links []string
	}

	type dataType struct {
		Header Header
		Footer Footer
		Body   Body
	}

	data := dataType{
		Header: Header{
			Title: "Test Title",
		},
		Body: Body{
			Title: "Body Test Title",
			Text:  "Text to be inject in the template file...",
		},
		Footer: Footer{
			Links: []string{"A", "B", "C"},
		},
	}

	err = tpl.ExecuteTemplate(os.Stdout, "composition.gohtml", data)
	if err != nil {
		log.Fatal(err)
	}
}
