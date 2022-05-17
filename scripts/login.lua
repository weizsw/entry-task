-- wrk.method = "POST"
-- wrk.body   = "username=test1&password=test"
-- wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
-- function request()

--     return wrk.format('POST', nil, wrk.headers, wrk.body)
  
-- end
function randomUser()
    suffix = math.random(1, 200)
    return "test" .. suffix
end

request = function() 
    user = randomUser()
    body = string.format('username=%s&password=test', user)
    headers = {}
    headers["Content-Type"] = "application/x-www-form-urlencoded"
    return wrk.format("POST", nil, headers, body)
end