Feature: Post Collection
    Scenario: POST /collections
        Given there are no collections
        When I POST "/collections"
            """
            {
                "name": "Coronavirus key indicators",
                "publish_date": "2020-05-05T14:58:29.317Z"
            }
            """
        Then the HTTP status code should be "201"

    Scenario: POST /collections where the collection name already exists
        Given I have these collections:
            """
            [
                {
                    "id": "abc123",
                    "name": "Coronavirus key indicators",
                    "publish_date": "2020-05-10T14:58:29.317Z"
                }
            ]
            """
        When I POST "/collections"
            """
            {
                "name": "Coronavirus key indicators",
                "publish_date": "2020-05-05T14:58:29.317Z"
            }
            """
        Then the HTTP status code should be "409"