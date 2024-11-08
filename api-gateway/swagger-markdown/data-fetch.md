## Store data from Indego

An endpoints which downloads fresh data from [Indego GeoJSON station status API](https://www.rideindego.com/stations/json/) and stores it inside PostgreSQL.

```bash
# this endpoint will be trigger every hour to fetch the data and insert it in the PostgreSQL database
POST http://localhost:3000/api/v1/indego-data-fetch-and-store-it-db
```

### Token 
Add HTTP header with Authorization 
```go

  // example:
	headers := map[string]string{
		"Authorization": "Bearer secret_token_static",
	},
```

### Response Code
| HTTP | Description                            |
|------|----------------------------------------|
| 200  | Fetch and store to database is success |
| 401  | Bad Authorization. Check token         |
