# Go path to wiki

Project that finds the shortest link between two Wikipedia pages. The program uses an wide search algorithm to find the fastest path, from the staring Wikipedia page to the ending page.

## To use

Create a bot account to mediawiki and uncomment the const from ./main.go and insert your own credentials

Then you can build the application with the command `go build`

Then the program can be be executed while giving it two parameters that are the start and the end

Example: `.\go-path-to-wiki.exe "Finland" "Apple"` Please note that the starting and ending pages have to be working Wikipedia pages

### Used resources

- https://pkg.go.dev/cgt.name/pkg/go-mwclient

- https://github.com/gosuri/uilive
