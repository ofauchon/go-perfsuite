package core

import (
   // "encoding/json"
  //  "fmt"
    "net/http"
    "log"
    //"core"

    //"github.com/gorilla/mux"
)


// Stuff needed th set up a custom handler, with context 


// The appContext will contains all the global objects 
// needed for processing requests (resources, loggers .... ) 
type appContext struct {
    controller *ControllerInstance
}

// This custom appHandler will contain :
// - an embedded field of appContext
// - a custom function that looks like original handler, but having *appContext in complement 
type appHandler struct {
    *appContext
    H func (*appContext, http.ResponseWriter, *http.Request) (int, error)
}

// Our custop appHandler
func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Updated to pass ah.appContext as a parameter to our handler type.
    status, err := ah.H(ah.appContext, w, r)    
    if err != nil {
        log.Printf("HTTP %d: %q", status, err)
        switch status {
        case http.StatusNotFound:
            http.NotFound(w, r)
            // And if we wanted a friendlier error page, we can
            // now leverage our context instance - e.g.
            // err := ah.renderTemplate(w, "http_404.tmpl", nil)
        case http.StatusInternalServerError:
            http.Error(w, http.StatusText(status), status)
        default:
            http.Error(w, http.StatusText(status), status)
        }
    }
}


// Our handlers


/*

func IndexHandler (w http.ResponseWriter, r *http.Request)  {
  //  fmt.Fprintln(w, "Welcome to the index %d", ah.controller)

}
*/
