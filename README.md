# Gmail Text Notifications
Send an SMS message when there are emails matching a given keyword(s). Uses the Twilio and Gmail APIs.

## Background
You know how some sites will let you sign up to be notified via email when something is back in stock? That's a great feature.

But when demand for an item greatly outpaces supply, the item is sure to go out of stock again shortly after the "back in stock" email is sent to those on the waitlist. This is especially true for limited edition or one-time items. This means that unless you see the email shortly after it's sent, you're likely to miss out on the item.

I don't particularly enjoy checking my email and only do so a couple times a day. I also don't have push notifications turned on for email. I am far, far more likely to see a text message, which I use pretty heavily. So I wrote this to send me a text message when certain emails come in.

But really it seemed like a good excuse to try a new programming language. I went with Go.

## Setup
Much of this makes use of the Gmail APIs, which require authentication. The easiest way to get set up to use this API is by going through the [Go Quickstart](https://developers.google.com/gmail/api/quickstart/go), which will have you download the necessary `credentials.json` (which I haven't uploaded to this repo for obvious reasons).

You'll also need a Twilio account and phone number. Using [my referral link](www.twilio.com/referral/XCX3Mu) will get each of us $10.

Once you've signed up for a Twilio account, you'll to store your Twilio credentials in a `config.json` - again, which I haven't uploaded to this repo for obvious reasons. Your `config.json` should look like this:
```
{
    "Twilio": {
        "AccountSID": "{{ YOUR_TWILIO_ACCOUNT_SID }}", "AuthToken": "{{ YOUR_TWILIO_AUTH_TOKEN }}", "PhoneNumber": "{{ YOUR_TWILIO_PHONE_NUMBER }}",
        "BaseURL": "https://api.twilio.com/2010-04-01"
    }
}
```

## Usage
```
$ go build main.go
$ ./main.go -q foo -phone +13125555555
```