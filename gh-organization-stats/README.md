# GitHub Organization Stats API
#### A simple API to get the stats for all the repositories in an organization in a single request, caching the data in a MongoDB database. This API will be used to showcase the club-wide stats of the current semester on [ufosc.org](https://ufosc.org).

## Install
Clone the repository (requires [git](https://git-scm.com/)):
```
https://github.com/ufosc/OpenWebServices
```

Navigate to the project directory and install the project dependencies (requires [Node.js](https://nodejs.org/en)):
```
cd OpenWebServices/gh-organization-stats
npm install
```
## Usage
<b>Starting the development server:</b>
```
npm run dev
```
You may access the server at http://localhost:8000

<b>Connect to MongoDB database by including this in .env file:</b>
```
DB_KEY=your_mongodb_connection_string
```

## API Endpoints
### POST ```/api/stats```
- Description: Calls the GitHub API to get the number of commits and pull requests for all the repositories in request body, aggregating the totals and storing them in the MongoDB database. Should be executed daily with a cron job.
- Request Body:
    ```
    {
      "startDate": "2024-09-16T00:00:00Z",
      "repos": [
        ["ufosc", "Alarm-Clock"],
        ["ufosc", "Club_Website_2"],
        ["ufosc", "OpenWebServices"]
      ]
    }
    ```
- Response:
  * ```200 OK```: Returns the aggregated stats for all the repositories in the request body.
      ```
      {
        "message": "Club-wide stats recorded successfully",
        "stats": {
          "totalCommits": 64,
          "totalOpenedPRs": 44
        }
      }
      ```
  * ```400 Bad Request```: Returns an error message if the request body is invalid.
      ```
      {
        "error": "Please provide startDate and repos."
      }
      ```
  * ```404 Not Found```: Returns an error message if one of the repositories in the request body are not found.
      ```
      {
        "error": "Error fetching stats for ufosc/fake-repo: Repository ufosc/fake-repo not found."
      }
      ```
  * ```500 Internal Server Error```: Returns an error message if there is an issue saving to the MongoDB collection or with connecting to the GitHub API.

### GET ```/api/stats```
- Description: Retrieves all the stats stored in the MongoDB database. Each semester's stats are stored in a separate document.
- Response:
  * ```200 OK```: Returns all the stats stored in the MongoDB database.
      ```
      {
        "stats": [
          {
            "_id": "6705e2d4c2ccbca1d8b52380",
            "start_date": "2024-09-16T00:00:00.000Z",
            "totalCommits": 9000,
            "totalOpenedPRs": 9000,
            "repos": [
              [
                "ufosc",
                "Alarm-Clock"
              ],
              [
                "ufosc",
                "Club_Website_2"
              ],
              ...
            ],
            "collection_date": "2024-10-09T00:00:00.000Z",
            "__v": 0
          }
        ]
      }
      ```
