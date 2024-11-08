
CREATE TABLE rideindego (
    last_update TIMESTAMP WITH TIME ZONE PRIMARY KEY,
    raw_data JSONB
); 

CREATE TABLE openweather (
    last_update TIMESTAMP WITH TIME ZONE PRIMARY KEY,
    raw_data JSONB
)
