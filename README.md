# HyperPlannerCourses

Experimenting with technologies meant to help query courses for HyperPlanner

- Golang
- MongoDB
- REST/AWS API Gateway/AWS Lambda
- GraphQL/AWS AppSync/AWS Lambda

### REST API Endpoints (requires api token)

Invoke URL: https://sq3oxmnhjc.execute-api.us-east-1.amazonaws.com/Development

Be sure to include `x-api-key` in the request header.

#### GET:
By default, hitting the endpoint returns all courses in the MongoDB Courses collection, as well as the number of courses. The shape of the json response is:
```
{
  "courses": [],
  "length": 0
}
```

Possible query string parameters:
- campus
- code
- credits
- title

(query example: https://sq3oxmnhjc.execute-api.us-east-1.amazonaws.com/Development?code=CSCI134)

#### POST:
Hitting the endpoint inserts all specified courses into the MongoDB Courses collection. To properly insert, the request body should have the following shape:
```
{
  "courses": [
    {
      "campus": "hmc",
      "code": "CSCI121",
      "credits": 3.0,
      "title": "Software Development"
    },
    {
      "campus": "hmc",
      "code": "CSCI134",
      "credits": 3.0,
      "title": "Operating Systems"
    }
  ]
}
```

The response will return the number of courses inserted in the following shape:
```
{
  "length": 0
}
```

### GraphQL API Endpoints (requires api token)

### Build process for AWS Lambda

To build `file.go` for AWS Lambda, we have to build for Linux since the binary for Lambda runs on Amazon Linux. If building on Windows, there are some permission issues that arise, but we can run a python script to fix this. Follow the steps below to build the binary and deploy to AWS Lambda. I use the UI to upload the zip file to Lambda.

First, run
```
$ GOARCH=amd64 GOOS=linux go build file.go
```
to build the binary.

Next, zip the binary so that you have a zip file called `file.zip`.

Then, run the `changePerms` function in `scripts/set-executable.py` to change the executable permissions necessary for Lambda. The `src` argument is `file.zip` and the `dst` argument is your output zip file name containing the executable `file` with the correct permissions.

Finally, upload the zip file to Lambda (either through UI or CLI).

