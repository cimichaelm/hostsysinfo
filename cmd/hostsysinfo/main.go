//
// File: main.go
//
//

package main

import (
)

//
// main entrypoint
//
func main() {

	obj := defaults()
     
    argparse(obj)
    
    app_begin(obj)
    app_run(obj)
    app_end(obj)
}
