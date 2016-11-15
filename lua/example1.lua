
-- b64 encoding
local function _b64enc( data )
    -- character table string
  local b='ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'

    return ( (data:gsub( '.', function( x ) 
        local r,b='', x:byte()
        for i=8,1,-1 do r=r .. ( b % 2 ^ i - b % 2 ^ ( i - 1 ) > 0 and '1' or '0' ) end
        return r;
    end ) ..'0000' ):gsub( '%d%d%d?%d?%d?%d?', function( x )
        if ( #x < 6 ) then return '' end
        local c = 0
        for i = 1, 6 do c = c + ( x:sub( i, i ) == '1' and 2 ^ ( 6 - i ) or 0 ) end
        return b:sub( c+1, c+1 )
    end) .. ( { '', '==', '=' } )[ #data %3 + 1] )
end


-- authentication header creation
local function _createBasicAuthHeader( username, password )
  -- the header format is "Basic <base64 encoded username:password>"
  local header = "Basic "
  local authDetails = _b64enc( username .. ":" .. password )
  header = header .. authDetails
  return header
end

local function urlEncode( str )
   if ( str ) then
      str = string.gsub( str, "\n", "\r\n" )
      str = string.gsub( str, "([^%w ])",
         function (c) return string.format( "%%%02X", string.byte(c) ) end )
      str = string.gsub( str, " ", "+" )
   end
   return str
end

local function pause(seconds)
  sl=math.random(1,seconds * 1000) 
  ---print("Thread ID ", k_GetId(), " sleep :", sl  )
  k_Sleep(sl)
end

--------------------------------


function rinit()
  http = require("http")
  re = require("re")
  math = require("math")
end


function rrun()

  while 1 do

    pause(5)
    ---print(k_GetId() .. "Resume from pause, starting action1")
    k_TransactionStart("action1_login_page")

    authent = _createBasicAuthHeader("olivier", "ah8dk2v6")
    resp, error_message = http.request("GET", "https://secure.aixmarseille.com/mail/", {
    headers={ 
        ["Accept"]="*/*", 
        ["Authorization"]=authent }
    }) 

    if (resp ~= nil ) then
      head_cook = resp.headers["Set-Cookie"] 
      if (head_cook ~= nil ) then 
        cook =re.match(resp.headers["Set-Cookie"] , "(roundcube_sessid=\\w+);")
        tok = re.match(resp.body , "name=\"_token\" value=\"(\\w+)\"")
        ---print(k_GetId(), "action1 OK StatusCode: ", resp.status_code," Cookie:", cook, "Tok:", tok,"\n")
        k_TransactionStop("action1_login_page", 1)
      else 
        print(k_GetId(), "action1 ERROR 1\n")
        k_TransactionStop("action1_login_page", 0)
      end
    else
        print(k_GetId(), "ERROR 2" .. error_message .. "\n")
      k_TransactionStop("action1_login_page", 0)
    end


    -- Step 2 : Authenticate
    if (tok ~= nil ) 
      then
      k_TransactionStart("action2_authenticate")
      user="olivier"
      password="secret" 
      url="https://secure.aixmarseille.com/mail/?_task=login"
      post="_token" .. tok .. "&_task=login&_action=login&_url=&_user=" .. user .."&_pass=" .. password
      resp, error_message = http.request("post", url, {
        headers={ 
         ["Accept"]="*/*", 
         ["Authorization"]=head_auth, 
         ["Cookie"]=cook,
         ["Content-Type"]="application/x-www-form-urlencoded"
        },
        form=post
      })
      k_TransactionStop("action2_authenticate", 1 )
    end
  
  end --while

---   for k, v in pairs( resp.headers ) do
---       print(">",k, v)
---  end

end --func


function rstop()
end


