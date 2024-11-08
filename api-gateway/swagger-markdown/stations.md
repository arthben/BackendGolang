## Snapshot of all stations at a specified time

Data for all stations as of 11am Universal Coordinated Time on September 1st, 2019:

```bash
GET http://localhost:3000/api/v1/stations?at=2019-09-01T10:00:00Z
```

This endpoint should respond as follows, with the actual time of the first snapshot of data on or after the requested time and the data:

```javascript
{
  at: '2019-09-01T10:00:00Z',
  stations: { /* As per the Indego API */ },
  weather: { /* As per the Open Weather Map API response for Philadelphia */ }
}
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
| 404  | Data Not Found                         |
