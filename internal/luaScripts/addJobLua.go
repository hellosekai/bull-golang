/**
 * @Description:
 * @FilePath: /bull-golang/internel/luaScripts/addJobLua.go
 * @Author: liyibing liyibing@lixiang.com
 * @Date: 2023-07-27 20:57:57
 */
package luaScripts

const AddJobLua = `
local jobId
local jobIdKey
local rcall = redis.call

local jobCounter = rcall("INCR", KEYS[4])

if ARGV[2] == "" then
  jobId = jobCounter
  jobIdKey = ARGV[1] .. jobId
else
  jobId = ARGV[2]
  jobIdKey = ARGV[1] .. jobId
  if rcall("EXISTS", jobIdKey) == 1 then
    return jobId .. "" -- convert to string
  end
end

-- Store the job.
rcall("HMSET", jobIdKey, "name", ARGV[3], "data", ARGV[4], "opts", ARGV[5], "timestamp", ARGV[6], "delay", ARGV[7], "priority", ARGV[9])

-- Check if job is delayed
local delayedTimestamp = tonumber(ARGV[8])
if(delayedTimestamp ~= 0) then
  local timestamp = delayedTimestamp * 0x1000 + bit.band(jobCounter, 0xfff)
  rcall("ZADD", KEYS[5], timestamp, jobId)
  rcall("PUBLISH", KEYS[5], delayedTimestamp)
else
  local target

  -- Whe check for the meta-paused key to decide if we are paused or not
  -- (since an empty list and !EXISTS are not really the same)
  local paused
  if rcall("EXISTS", KEYS[3]) ~= 1 then
    target = KEYS[1]
    paused = false
  else
    target = KEYS[2]
    paused = true
  end

  -- Standard or priority add
  local priority = tonumber(ARGV[9])
  if priority == 0 then
      -- LIFO or FIFO
    rcall(ARGV[10], target, jobId)
  else
    -- Priority add
    rcall("ZADD", KEYS[6], priority, jobId)
    local count = rcall("ZCOUNT", KEYS[6], 0, priority)

    local len = rcall("LLEN", target)
    local id = rcall("LINDEX", target, len - (count-1))
    if id then
      rcall("LINSERT", target, "BEFORE", id, jobId)
    else
      rcall("RPUSH", target, jobId)
    end

  end

  -- Emit waiting event (wait..ing@token)
  rcall("PUBLISH", KEYS[1] .. "ing@" .. ARGV[11], jobId)
end

return jobId .. "" -- convert to string

`
