local key = KEYS[1]
local cntKey = key.."cnt"
local val = ARGV[1]
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    return -2
elseif ttl == -2 or ttl < 540 then
    redis.call("set",key,val)
    redis.call("expire",key,600)
    redis.call("set",cntKey,3)
    redis.call("expire",cntKey,600)
    return 0
else
    return -1
end