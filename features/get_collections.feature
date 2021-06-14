Feature: Get Collections
    Scenario: GET /collections
        Given I have these collections:
            """
            [
                {
                    "id": "abc123",
                    "name": "LMSV1",
                    "publish_date": "2020-05-10T14:58:29.317Z"
                },
                {
                    "id": "abc124",
                    "name": "LMSV2",
                    "publish_date": "2020-05-05T14:58:29.317Z"
                },
                {
                    "id": "abc125",
                    "name": "LMSV3",
                    "publish_date": "2020-05-08T14:58:29.317Z"
                }
            ]
            """
        When I GET "/collections"
        Then the HTTP status code should be "200"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following JSON response:
            """
            {
                "count": 3,
                "limit": 20,
                "offset": 0,
                "total_count": 3,
                "items": [
                    { "id": "abc123", "name": "LMSV1", "publish_date": "2020-05-10T14:58:29.317Z" },
                    { "id": "abc124", "name": "LMSV2", "publish_date": "2020-05-05T14:58:29.317Z" },
                    { "id": "abc125", "name": "LMSV3", "publish_date": "2020-05-08T14:58:29.317Z" }
                ]
            }
            """
