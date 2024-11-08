## Snapshot of one station at a specific time

Data for a specific station (by its `kioskId`) at a specific time:

```bash
GET http://localhost:3000/api/v1/stations/{kioskId}?at=2019-09-01T10:00:00Z
```

The response should be the first available on or after the given time, and should look like:

```javascript
{
  at: '2019-09-01T10:00:00',
  station: { /* Data just for this one station as per the Indego API */ },
  weather: { /* As per the Open Weather Map API response for Philadelphia */ }
}
```

Include an `at` property in the same format indicating the actual time of the snapshot.

If no suitable data is available a 404 status code should be given.


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
