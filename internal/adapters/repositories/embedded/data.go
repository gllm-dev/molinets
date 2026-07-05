package embedded

import _ "embed"

//go:embed data/grinders.json
var GrindersJSON []byte

//go:embed data/methods.json
var MethodsJSON []byte
