package.path = package.path .. "/scripts"

local JSON = require("JSON")
math.randomseed(os.time())

local gateway = "BBBB"

-- according to counts specified in the seeder

local max_user_suffix = 49999
local max_city_suffix = 24
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
        FunctionName = "login",
        Input = {
            Email = email_prefix .. tostring(id),
            Password = "Password" .. tostring(id)
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
            FunctionName = "recommendHotelsRate",
            Input = {
                City = city_prefix .. tostring(city_id),
                Count =  6
            }
        }
    else
        path = gateway
        param = {
            FunctionName = "recommendHotelsLocation",
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

local function search_hotel()
    local city_id = math.random(0, max_city_suffix)
    local method = "GET"
    local param = {
        FunctionName = "getHotelsInCity",
        Input = city_prefix .. tostring(city_id)
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
        FunctionName = "reserveRoom",
        Input = {
            DateIn = date1,
            DateOut = date2,
            UserEmail = email_prefix .. tostring(email_id),
            HotelName = hotel_prefix .. tostring(hotel_id),
            CityName = city_prefix .. tostring(city_id),
            RoomId = "Room" .. tostring(room_id)
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
        FunctionName = "registerUser",
        Input = {
            FirstName = "NewFirstName",
            LastName = "NewLastName",
            Email = registeredUserEmail,
            Password = registeredUserPassword
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
        FunctionName = "deleteUser",
        Input = {
            Email = registeredUserEmail,
            Password = registeredUserPassword
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
        FunctionName = "setHotelRate",
        Input = {
            Rate = city_id % 6,
            CityName = city_prefix .. tostring(city_id),
            HotelName = hotel_prefix .. tostring(hotel_id)
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
        FunctionName = "getUserReservations",
        Input = email_prefix .. tostring(user_id)
    }
    local body = JSON:encode(param)
    local headers = {}
    headers["Content-Type"] = "application/json"

    return wrk.format(method, gateway, headers, body)
end

---@diagnostic disable-next-line: lowercase-global
request = function()
    local req = Id
    Id = (Id + 1) % function_n

    if req == 0 then
        return login()
    elseif req == 1 then
        return recommend()
    elseif req == 2 then
        return search_hotel()
    elseif req == 3 then
        return reserve()
    elseif req == 4 then
        return add_user()
    elseif req == 5 then
        return delete_user()
    elseif req == 6 then
        return set_hotel_rate()
    elseif req == 7 then
        return get_user_reservations()
    end
end

-- ---@diagnostic disable-next-line: lowercase-global
-- response = function(code, header, body)
--     print(code)
--     print(body)
-- end
