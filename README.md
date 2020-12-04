# Gmail Text Notifications
Send an SMS message when there are emails matching a given keyword(s). Uses the Twilio and Gmail APIs.

## Background
You know how some sites will let you sign up to be notified via email when something is back in stock? That's a great feature.

But when demand for an item greatly outpaces supply, the item is sure to go out of stock again shortly after the "back in stock" email is sent to those on the waitlist. This is especially true for limited edition or one-time items. This means that unless you see the email shortly after it's sent, you're likely to miss out on the item.

I don't particularly enjoy checking my email and only do so a couple times a day. I also don't have push notifications turned on for email. I am far, far more likely to see a text message, which I use pretty heavily. So I wrote this to send me a text message when certain emails come in.

But really it seemed like a good excuse to try a new programming language. I went with Go.

## Usage
```
$ go build main.go
$ ./main.go -q foo -phone +13125555555
```