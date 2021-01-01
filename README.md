# Discord Webhook Proxy
## Okay, so here's the thing:

Many websites / services offer the ability to call webhooks on certain actions.  
Unfortunately, **only a few of them offer webhook support for Discord**, so for most websites / services you have to write a small tool to handle these requests.

For example, I'm using a **Raspberry Pi** with [MotionEyeOS](https://github.com/ccrisan/motioneyeos) as a **surveillance camera**, and I wanted to use Discord webhooks to notify me when motion was detected, so I had to develop a [small tool](https://github.com/darmiel/gomera).

In the middle of this, however, it occurred to me that this isn't the first time I've had to go through such a roundabout way just to receive a simple webhook. So I thought to write a "*universal*" Discord webhook "*proxy*", which can forward any kind of request to a valid Discord webhook.

For this, the JSON of the webhook is stored with placeholders, like {{ test }}, in a sqlite3 database, which will be replaced by the URL query parameters later and sent to the specified webhook:

### Example webhook data
webhook1:
```json
{
  "content": "Hey @everyone, this is a test! #{{ test }}",
  "username": "Webhook",
  "embeds": [
    {
      "title": "View test online!",
      "description": "This is a description for test: {{ test }}",
      "url": "https://test.com",
      "color": 3971831,
      "fields": [
        {
          "name": "Test",
          "value": "{{ test }}",
          "inline": true
        }
      ],
      "author": {
        "name": "Webhook",
        "icon_url": "https://image.com/test.png"
      }
    }
  ]
}
```
If you now send a request to:  
https://my-service.com/webhook/webhook1/{webhook_url}/?test=abc  
It will replace the placeholders "test" `{{ test }}` with "abc" and then send the webhook `{webhook_url}`.

The placeholders will be determined by **URL query parameters**, or by **POST application/json, application/x-www-form-urlencoded**