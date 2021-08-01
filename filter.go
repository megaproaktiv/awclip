package awclip

var discriminatedCommandsList = []string{
	"generate-credential-report",

}

// DiscriminatedCommand
// check whether a command is cachable
func DiscriminatedCommand( command *string) bool{
	return contains(discriminatedCommandsList, *command)
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}