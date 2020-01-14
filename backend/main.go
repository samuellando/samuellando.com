// main.go

package main


func main() {
  a := App{}
  a.Initialize(
    "pageAdmin",
    "",
    "pages")

    a.Run(":8080")
}
