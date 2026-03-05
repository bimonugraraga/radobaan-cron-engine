package cronengine

import "github.com/redis/go-redis/v9"

var AcquireScript = redis.NewScript(`
local prev = redis.call("SETBIT", KEYS[1], ARGV[1], 1)

if prev == 1 then
    return 0
end

if tonumber(ARGV[2]) > 0 then
    local ttl = redis.call("TTL", KEYS[1])
    if ttl < 0 then
        redis.call("EXPIRE", KEYS[1], ARGV[2])
    end
end

return 1
`)
