Feature: Get Collection events
    Scenario: GET /collections/{collection_id}/events
        Given I have a collection with ID "coronaviruskeyindicators-5d57ce55" with the following events:
            """
            [
                {
                    "date": "2020-05-05T14:58:29.317Z",
                    "type": "CREATED",
                    "email": "person.name@ons.gov.uk"
                }
            ]
            """
        When I GET "/collections/coronaviruskeyindicators-5d57ce55/events"
        Then the HTTP status code should be "200"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following JSON response:
            """
            {
                "count": 1,
                "limit": 20,
                "offset": 0,
                "total_count": 1,
                "items": [
                    {
                        "date": "2020-05-05T14:58:29.317Z",
                        "type": "CREATED",
                        "email": "person.name@ons.gov.uk"
                    }
                ]
            }
            """
