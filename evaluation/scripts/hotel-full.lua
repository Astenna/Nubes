---@diagnostic disable: lowercase-global
package.path = package.path .. "/scripts"

local JSON = require("JSON")
local base64 = require("base64")
math.randomseed(os.time())

local gateway = "https://qzod5szgz6hgwox3suztsgc4ze0mkefm.lambda-url.us-east-1.on.aws/"

-- according to counts specified in the seeder

local max_user_suffix = 49999
local max_city_suffix = 4
local max_hotel_suffix = 99
local max_room_suffix = 24
local email_prefix = "Email"
local city_prefix = "Milano"
local hotel_prefix = "Bruschetti"

local registeredUserEmail = "registeredUserEmail" .. tostring(math.random(10))
local registeredUserPassword = "registeredUserPassword" .. tostring(math.random(10))

local function_n = 8

Id = math.random(function_n)

local function login()
    local id = math.random(0, max_user_suffix)
    local method = "GET"
    local param = {
        FunctionName = "UserVerifyPassword",
        Input = {
            Id = email_prefix .. tostring(id),
            Parameter = "Password" .. tostring(id)
        }
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
            FunctionName = "CityGetHotelsWithBestRates",
            Input = {
                Id = city_prefix .. tostring(city_id),
                Parameter = 6
            }
        }
    else
        path = gateway
        param = {
            FunctionName = "CityGetHotelsCloseTo",
            Input = {
                Id = city_prefix .. tostring(city_id),
                Parameter = {
                    Count = 6,
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

local function search_hotel()
    local city_id = math.random(0, max_city_suffix)
    local method = "GET"
    local param = {
        FunctionName = "ReferenceGetStubs",
        Input = {
            OwnerId = city_prefix .. tostring(city_id),
            OwnerTypeName = "City",
            OtherTypeName = "Hotel",
            ReferringFieldName = "City",
            IsManyToMany = false
        }
    }
    local body = JSON:encode(param)
    local headers = {}
    headers["Content-Type"] = "application/json"

    return wrk.format(method, gateway, headers, body)
end

local function get_two_consecutive_days_in_year(year, length)
    local day1 = math.random(1, 28)
    local month1 = math.random(1, 12)
    local month2 = month1
    local day2 = day1 + length
    local year2 = year

    -- handle the cases when day2 is out of range of a month
    if day2 > 30 and (month1 == 4 or month1 == 6 or month1 == 9 or month1 == 11) then
        month2 = month2 + 1
        day2 = day2 - 30
        if month2 > 12 then
            month2 = 1
            year2 = year + 1
        end
    elseif day2 > 28 and month1 == 2 then
        if (year % 400 == 0) or ((year % 4 == 0) and (year % 100 ~= 0)) then
            -- leap year
            if day2 > 29 then
                day2 = day2 - 29
                month2 = month2 + 1
                if month2 > 12 then
                    month2 = 1
                    year2 = year + 1
                end
            end
        else
            day2 = day2 - 28
            month2 = month2 + 1
            if month2 > 12 then
                month2 = 1
                year2 = year + 1
            end
        end
    elseif day2 > 31 then
        day2 = day2 - 31
        month2 = month2 + 1
        if month2 > 12 then
            month2 = 1
            year2 = year + 1
        end
    end

    -- add leading zeros
    local day1_str = tostring(day1)
    local day2_str = tostring(day2)
    if day1 < 10 then
        day1_str = "0" .. day1_str
    end
    if day2 < 10 then
        day2_str = "0" .. day2_str
    end

    local month1_str = tostring(month1)
    local month2_str = tostring(month2)
    if month1 < 10 then
        month1_str = "0" .. month1_str
    end
    if month2 < 10 then
        month2_str = "0" .. month2_str
    end

    return tostring(year) .. "-" .. month1_str .. "-" .. day1_str,
        tostring(year2) .. "-" .. month2_str .. "-" .. day2_str
end

local function reserve()
    local email_id = math.random(0, max_user_suffix)
    local city_id = math.random(0, max_city_suffix)
    local hotel_id = math.random(0, max_hotel_suffix)
    local room_id = math.random(0, max_room_suffix)

    -- in 50% of cases try to reserve a room in dates
    -- where the room is likely to be fully booked
    local coin = math.random()
    local date1, date2 = "", ""
    if coin < 0.5 then
        date1, date2 = get_two_consecutive_days_in_year(2023, math.random(1, 14))
    else
        date1, date2 = get_two_consecutive_days_in_year(2024, math.random(1, 14))
    end

    local method = "GET"
    local param = {
        FunctionName = "Export",
        Input = {
            TypeName = "Reservation",
            Parameter = {
                DateIn = date1,
                DateOut = date2,
                User = email_prefix .. tostring(email_id),
                RoomId = city_prefix ..
                    tostring(city_id) .. "_" .. hotel_prefix .. tostring(hotel_id) .. "_" .. "Room" .. tostring(room_id)
            }
        }
    }

    local body = JSON:encode(param)
    local headers = {}
    headers["Content-Type"] = "application/json"
    return wrk.format(method, gateway, headers, body)
end

local function add_user()
    local method = "GET"
    local param = {
        FunctionName = "Export",
        Input = {
            TypeName = "User",
            Parameter = {
                FirstName = "NewFirstName",
                LastName = "NewLastName",
                Email = registeredUserEmail,
                Password = registeredUserPassword
            }
        }
    }
    local body = JSON:encode(param)
    local headers = {}
    headers["Content-Type"] = "application/json"

    return wrk.format(method, gateway, headers, body)
end


local function delete_user()
    local method = "GET"
    local param = {
        FunctionName = "Delete",
        Input = {
            TypeName = "User",
            Parameter = {
                Email = registeredUserEmail,
                Password = registeredUserPassword
            }
        }
    }
    local body = JSON:encode(param)
    local headers = {}
    headers["Content-Type"] = "application/json"

    return wrk.format(method, gateway, headers, body)
end

local function set_hotel_rate()
    local city_id = math.random(0, max_city_suffix)
    local hotel_id = math.random(0, max_hotel_suffix)
    local method = "GET"
    local param = {
        FunctionName = "SetField",
        Input = {
            Id = city_prefix .. tostring(city_id) .. "_" .. hotel_prefix .. tostring(hotel_id),
            FieldName = "Rate",
            TypeName = "Hotel",
            Value = city_id % 6
        }
    }
    local body = JSON:encode(param)
    local headers = {}
    headers["Content-Type"] = "application/json"

    return wrk.format(method, gateway, headers, body)
end

local function get_user_reservations()
    local user_id = math.random(0, max_user_suffix)
    local method = "GET"
    local param = {
        FunctionName = "ReferenceGetStubs",
        Input = {
            OwnerId = email_prefix .. tostring(user_id),
            OwnerTypeName = "User",
            OtherTypeName = "Reservation",
            ReferringFieldName = "Reservations",
            IsManyToMany = true
        }
    }
    local body = JSON:encode(param)
    local headers = {}
    headers["Content-Type"] = "application/json"

    return wrk.format(method, gateway, headers, body)
end

Req = "none"

---@diagnostic disable-next-line: lowercase-global
request = function()
    local req = Id
    Id = (Id + 1) % function_n

    if req == 0 then
        Req = "login"
        return login()
    elseif req == 1 then
        Req = "recommend"
        return recommend()
    elseif req == 2 then
        Req = "search_hotel"
        return search_hotel()
    elseif req == 3 then
        Req = "reserve"
        return reserve()
    elseif req == 4 then
        Req = "add_user"
        return add_user()
    elseif req == 5 then
        Req = "delete_user"
        return delete_user()
    elseif req == 6 then
        Req = "set_hotel_rate"
        return set_hotel_rate()
    elseif req == 7 then
        Req = "get_user_reservations"
        return get_user_reservations()
    end
end

function printTable(t)
    for k, v in pairs(t) do
        if type(v) == "table" then
            print(k .. ":")
            printTable(v)  -- Recursive call for nested tables
        else
            print(k .. ": " .. tostring(v))
        end
    end
end


-- ---@diagnostic disable-next-line: lowercase-global
-- response = function(code, header, body)
--     if code == 200 then
--         local data = JSON:decode(body)

--         if data.FunctionError ~= nil then
--             local payload = data.Payload
--             local decodedPayload = base64.decode(payload)
--             print(string.format("FAILED %s", Req))
--             -- printTable(header)
--             print(decodedPayload)
--         else
--             print(string.format("OK %s", Req))
--         end

--         -- print(data)
--     else
--         print(string.format("FAILED(%d) %s", code, Req))
--         print(body)
--     end
-- end
