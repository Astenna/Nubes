package.path = package.path .. "/scripts"

local JSON = require("JSON")
math.randomseed(os.time())

local gateway = "CCC"

-- according to counts specified in the seeder
local max_city_suffix = 24
local city_prefix = "Milano"

local function_n = 2

local function search_hotel()
    local city_id = math.random(0, max_city_suffix)
    local method = "GET"
    local param = {
        FunctionName = "getHotelsInCitySimple",
        Input = city_prefix .. tostring(city_id)
    }
    local body = JSON:encode(param)
    local headers = {}
    headers["Content-Type"] = "application/json"

    return wrk.format(method, gateway, headers, body)
end

local function recommend()
    local recommend_rate = math.random(0, 1)
    local city_id = math.random(0, max_city_suffix)
    local method = "GET"
    local path = ""
    local param = {}

    if recommend_rate == 0 then
        path = gateway
        param = {
            FunctionName = "recommendHotelsRateSimple",
            Input = {
                City = city_prefix .. tostring(city_id),
                Count = 6
            }
        }
    else
        path = gateway
        param = {
            FunctionName = "recommendHotelsLocationSimple",
            Input = {
                City = city_prefix .. tostring(city_id),
                Count = 6,
                Coordinates = {
                    Longitude = (-1) * math.random(0, 90) + math.random(0, 89) + math.random(),
                    Latitude = (-1) * math.random(0, 180) + math.random(0, 179) + math.random()
                }
            }
        }
    end

    local body = JSON:encode(param)
    local headers = {}
    headers["Content-Type"] = "application/json"

    return wrk.format(method, path, headers, body)
end

Id = math.random(function_n)

---@diagnostic disable-next-line: lowercase-global
request = function()
    Id = (Id + 1) % function_n
    local req = Id

    if req == 0 then
        return recommend()
    elseif req == 1 then
        return search_hotel()
    end
end

-- ---@diagnostic disable-next-line: lowercase-global
-- response = function(code, header, body)
--     print(code)
--     print(body)
-- end
