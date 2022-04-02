/*

Package goroutine discovers and returns information about either all goroutines
or the caller's goroutine. The information provided by the Goroutine type
consists of a unique ID, the state, the name of the topmost (most recent)
function in the call stack, as well as the stack backtrace. For goroutines other
than the main goroutine (the one with ID 1) the creating function as well as
location are additionally provided.

*/
package goroutine
