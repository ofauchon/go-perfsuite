function rinit()
  local http = require("http")
end

function rrun()
  k_CounterStart("load_http_google")
  local http = require("http")
  response, error_message = http.request("GET", "http://www.example.com", {
    query="q=test",
    headers={
        Accept="*/*"
    }   
  })  
 print("LEN: " ..  string.len(response.body))

  k_CounterEnd("load_http_google")
end 


function rstop()
end
