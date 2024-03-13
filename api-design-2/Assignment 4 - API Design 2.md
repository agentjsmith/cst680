# 1. API Changes

I based these changes mostly on the HAL standard which defines a format for hyperlinking in JSON documents, and the HAL-FORMS draft standard, which adds the ability to describe methods other than "GET".  I flat out violated RFC-8288 because I could not stand the sight of URLs as hash keys (for relationships not registered with IANA), and even the examples in the HAL document do not respect that limitation.  Unless a popular client comes along that requires that ~~bureaucratic nonsense~~ behavior, I am happy with my choice.  Even XML allowed giving a short name to a namespace URL.

To support this change, I modified the ID fields to use URLs in the API instead of integers.  The system can still use integers, GUIDs, or anything else internally.  The other related change I made is to assume that the API will automatically assign a suitable internal ID when an object is added and return the new URL to the client in response.

Here are hypothetical JSON documents produced by the modified API:
## Voter API
```json
<!-- GET https://voterserver/voters -->
{
	"_embedded": {
		"item": {"title": "Count Chocula", "href": "/voters/1"},
		"item": {"title": "Captain Crunch", "href": "/voters/2"}
	},
	"_templates": {
		"default": {
			"title": "Add Voter",
			"method": "POST",
			"contentType": "application/json",
			"properties": [
				{
					"name": "first_name",
					"required": true, 
					"prompt": "First Name"
				},
				{
					"name": "last_name",
					"required": true, 
					"prompt": "Last Name"
				},
			]
		}
	},
	"_links": {
		"polls": {
			"title": "Open Polls",
			"href": "https://pollerserver/polls"
		}
	}
}

<!-- GET https://voterserver/voters/1 -->
{
	"id": 1,
	"first_name": "Count",
	"last_name": "Chocula",
	"_links": {
		"self": {"href": "/voters/1"},
		"history": {"href": "/voters/1/polls"}
	},
	"_templates": {
		"delete": {
			"title": "Delete Voter 1",
			"method": "DELETE",
			"contentType": "application/json",
		},
		"default": {
			"title": "Update Voter 1",
			"method": "PUT",
			"contentType": "application/json",
			"properties": [
				{
					"name": "first_name",
					"required": true, 
					"prompt": "First Name"
				},
				{
					"name": "last_name",
					"required": true, 
					"prompt": "Last Name"
				},
			]
		}
	}
}

<!-- GET https://voterserver/voters/1/polls -->
{
	"_embedded": {
		"history": [
			{
				"item": {
					"date": "20230506", 
					"href": "https://pollserver/polls/1"
				}
			},
			{
				"item": {
					"date": "20240314", 
					"href": "https://pollserver/polls/2"
				}
			}
		]
	},
	"_templates": {
		"default": {
			"title": "Add Vote History to Voter 1",
			"method": "POST",
			"contentType": "application/json",
			"properties": [
				{
					"name": "poll",
					"required": true, 
					"prompt": "Poll URL"
				},
				{
					"name": "date",
					"required": true, 
					"prompt": "Date"
				},
			]
		}
	}
}
```

## Poll API
```json

<!-- GET https://pollserver/polls -->
{
	"_embedded": {
		"polls": [
			{"item": {"title": "Favorite Color", "href": "/polls/1"}},
			{"item": {"title": "Favorite Pet", "href": "/polls/2"}},
		]
	},
	"_templates": {
		"default": {
			"title": "Add Poll",
			"method": "POST",
			"contentType": "application/json",
			"properties": [
				{
					"name": "title",
					"required": true, 
					"prompt": "Poll Title"
				},
				{
					"name": "question",
					"required": true, 
					"prompt": "Question"
				},
			]
		}
	}
}

<!-- GET https://pollserver/polls/1 -->
{
	"id": 1,
	"title": "Favorite Color",
	"question": "What is your favorite color channel?",
	"_links": {
		"self": {"href": "/polls/1"}
	},
	"_embedded": {
		"options": [
			{
				"item": {
					"title": "Red",
					"href": "/polls/1/options/1"
				}
			},
			{
				"item": {
					"title": "Green",
					"href": "/polls/1/options/2"
				}
			},
			{
				"item": {
					"title": "Blue",
					"href": "/polls/1/options/3"
				}
			},
			{
				"item": {
					"title": "Alpha",
					"href": "/polls/1/options/4"
				}
			},
		]
	},
	"_templates": {
		"add_option": {
			"title": "Add Option to Poll 1",
			"method": "POST",
			"target": "/polls/1/options",
			"properties": [
				{
					"name": "text",
					"required": "true",
					"prompt": "Option Text"
				}
			]
		},
		"delete": {
			"title": "Delete Poll 1",
			"method": "DELETE",
			"contentType": "application/json",
		},
		"default": {
			"title": "Update Poll 1",
			"method": "PUT",
			"contentType": "application/json",
			"properties": [
				{
					"name": "title",
					"required": true, 
					"prompt": "Poll Title"
				},
				{
					"name": "question",
					"required": true, 
					"prompt": "Question"
				},
			]
		}
	}
}

<!-- GET https://pollserver/polls/1/options/1 -->
{
	"id": 1,
	"text": "Red",
	"_links": {
		"self": {"href": "/polls/1/options/1"}
	},
	"_templates": {
		"delete": {
			"title": "Delete Option 1",
			"method": "DELETE",
			"contentType": "application/json",
		},
		"default": {
			"title": "Update Option 1",
			"method": "PUT",
			"contentType": "application/json",
			"properties": [
				{
					"name": "text",
					"required": true, 
					"prompt": "Option Text"
				}
			]
		}
	}
}
```

## Vote API
```json
<!-- GET https://voteserver/votes -->

{
	"_embedded": {
		"votes": [
			{
				"item": {
					"poll": {"href": "https://pollserver/polls/1"},
					"option":
						{"href": "https://pollserver/polls/1/options/1"},
					"voter": {"href": "https://voterserver/voters/1"}
				}
			},
			{
				"item": {
					"poll": {"href": "https://pollserver/polls/2"},
					"option":
						{"href": "https://pollserver/polls/1/options/4"},
					"voter": {"href": "https://voterserver/voters/1"}
				}
			},
		]
	},
	"_templates": {
		"default": {
			"title": "Cast Vote",
			"method": "POST",
			"contentType": "application/json",
			"properties": [
				{
					"name": "voter",
					"required": true, 
					"prompt": "Voter URL"
				},
				{
					"name": "poll",
					"required": true, 
					"prompt": "Poll URL"
				},
				{
					"name": "option",
					"required": true, 
					"prompt": "Option URL"
				},
			]
		}
	},
	"_links": {
		"voters": {
			"title": "Voter Roll",
			"href": "https://voterserver/voters"
		},
		"polls": {
			"title": "Open Polls",
			"href": "https://pollerserver/polls"
		}
	}
}
```

# 2. How the APIs would work together better

Now a consumer of any of the APIs will get back information about how to use that API in a machine-usable format, but with human-readable labels.  So the person developing the API consumer only needs to know the endpoints for each.  In fact, the Votes API which I envision as the primary user-facing API, links to the other two APIs that it depends on, so a developer could start there and discover the entire API.

Better yet, the links in each JSON document can depend on the permissions of the  logged-in user.  So for an individual voter who could only update their own record, they would only see links for updating their own record.
# 3. Helpful Sources

+ HATEOAS.  https://en.wikipedia.org/wiki/HATEOAS
- JSON Hypertext Application Language.  https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-11
- Web Linking. https://www.rfc-editor.org/rfc/inline-errata/rfc8288.html
- ChatGPT 4. https://chat.openai.com/share/006310d8-59f7-41bc-8236-329a97c538d2
- The HAL-FORMS Media Type.  https://rwcbook.github.io/hal-forms/
- URI Template.  https://www.rfc-editor.org/rfc/rfc6570.html