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

    Scenario: GET /collections when no collections exist
            Given I have these collections:
                """
                []
                """
            When I GET "/collections"
            Then the HTTP status code should be "200"
            And the response header "Content-Type" should be "application/json; charset=utf-8"
            And I should receive the following JSON response:
                """
                {
                    "count": 0,
                    "limit": 20,
                    "offset": 0,
                    "total_count": 0,
                    "items": []
                }
                """

    Scenario: GET /collections with order
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
        When I GET "/collections?order_by=publish_date"
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
                    { "id": "abc124", "name": "LMSV2", "publish_date": "2020-05-05T14:58:29.317Z" },
                    { "id": "abc125", "name": "LMSV3", "publish_date": "2020-05-08T14:58:29.317Z" },
                    { "id": "abc123", "name": "LMSV1", "publish_date": "2020-05-10T14:58:29.317Z" }
                ]
            }
            """

    Scenario: GET /collections with invalid order
        When I GET "/collections?order_by=FUBAR"
        Then the HTTP status code should be "400"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following JSON response:
            """
            {
                "errors":[ {"message":  "invalid order_by"}]
            }
            """

    Scenario: GET /collections with name search
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
        When I GET "/collections?name=LMSV3"
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
                    { "id": "abc125", "name": "LMSV3", "publish_date": "2020-05-08T14:58:29.317Z" }
                ]
            }
            """

    Scenario: GET /collections with name search - multiple results
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
                    "name": "LMSV1 second edition",
                    "publish_date": "2020-05-05T14:58:29.317Z"
                },
                {
                    "id": "abc125",
                    "name": "LMSV2",
                    "publish_date": "2020-05-08T14:58:29.317Z"
                }
            ]
            """
        When I GET "/collections?name=LMSV1"
        Then the HTTP status code should be "200"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following JSON response:
            """
            {
                "count": 2,
                "limit": 20,
                "offset": 0,
                "total_count": 2,
                "items": [
                    { "id": "abc123", "name": "LMSV1", "publish_date": "2020-05-10T14:58:29.317Z" },
                    { "id": "abc124", "name": "LMSV1 second edition", "publish_date": "2020-05-05T14:58:29.317Z" }
                ]
            }
            """

    Scenario: GET /collections with a name search that's more than 64 characters long
        When I GET "/collections?name=0123456789012345678901234567890123456789012345678901234567890123456789"
        Then the HTTP status code should be "400"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following JSON response:
            """
            {
                "errors":[ {"message":  "name search text is >64 chars"}]
            }
            """

    Scenario: GET a specific collection
        Given I have these collections:
        """
            [
                {
                    "id": "00112233-4455-6677-8899-aabbccddeeff",
                    "name": "LMSV1",
                    "e_tag": "123",
                    "publish_date": "2020-05-10T14:58:29.317Z"
                },
                {
                    "id": "10112233-4455-6677-8899-aabbccddeeff",
                    "name": "LMSV2",
                    "e_tag": "567",
                    "publish_date": "2020-05-05T14:58:29.317Z"
                },
                {
                    "id": "20112233-4455-6677-8899-aabbccddeeff",
                    "name": "LMSV3",
                    "e_tag": "456",
                    "publish_date": "2020-05-08T14:58:29.317Z"
                }
            ]
        """
        When I GET "/collections/00112233-4455-6677-8899-aabbccddeeff"
        Then the HTTP status code should be "200"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And the response header "Etag" should be "123"
        And I should receive the following JSON response:
        """
        {
            "id": "00112233-4455-6677-8899-aabbccddeeff",
            "name": "LMSV1",
            "e_tag":"123",
            "publish_date": "2020-05-10T14:58:29.317Z"
        }
        """

    Scenario: GET a specific collection that does not exist
        Given I have these collections:
        """
            [
                {
                    "id": "10112233-4455-6677-8899-aabbccddeeff",
                    "name": "LMSV2",
                    "publish_date": "2020-05-05T14:58:29.317Z"
                },
                {
                    "id": "20112233-4455-6677-8899-aabbccddeeff",
                    "name": "LMSV3",
                    "publish_date": "2020-05-08T14:58:29.317Z"
                }
            ]
        """

        When I GET "/collections/00112233-4455-6677-8899-aabbccddeeff"
        Then the HTTP status code should be "404"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following JSON response:
        """
        {
            "errors":[ {"message":  "collection not found"}]
        }
        """
