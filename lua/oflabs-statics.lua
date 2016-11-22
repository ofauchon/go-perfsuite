local function pause(seconds)
  sl=math.random(1,seconds * 1000) 
  ---print("Thread ID ", k_GetId(), " sleep :", sl  )
  k_Sleep(sl)
end

function rinit()
  http = require("http")
  re = require("re")
  math = require("math")
end




function rrun()

  while 1 do

    pause(5)
 
    ---print(k_GetId() .. "Resume from pause, starting action1")
    k_TransactionStart("action1_home")
    resp, error_message = http.request("GET", "https://www.oflabs.com", {
    headers={ 
        ["Accept"]="*/*", 
        ["Authorization"]=authent }
    }) 
    k_TransactionStop("action1_home", resp ~= nil)

    stats = http.get_stats()
    for k, v in pairs( stats ) do
       print(">",k, v)
    end



    ---print(k_GetId() .. "Resume from pause, starting action1")
    k_TransactionStart("action1_image")
    resp, error_message = http.request("GET", "https://www.oflabs.com/content/images/2016/03/081209-light-bulb-02-2.jpg", {
    headers={ 
        ["Accept"]="*/*", 
        ["Authorization"]=authent }
    }) 
    k_TransactionStop("action1_image", resp ~= nil )

  
  end --while


end --func


function rstop()
end


